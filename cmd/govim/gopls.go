package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/govim/govim/cmd/govim/config"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/fakenet"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/jsonrpc2"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/protocol"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/span"
	"github.com/govim/govim/cmd/govim/internal/util"
)

func (s *goplsServer) start(initParams *protocol.ParamInitialize) error {
	logfile, err := s.g.createLogFile("gopls")
	if err != nil {
		return err
	}
	logfile.Close()
	s.g.Logf("gopls log file: %v", logfile.Name())

	s.g.ChannelExf("let s:gopls_logfile=%q", logfile.Name())

	goplsArgs := []string{"-rpc.trace", "-logfile", logfile.Name()}
	if flags, err := util.Split(os.Getenv(string(config.EnvVarGoplsFlags))); err != nil {
		s.g.Logf("invalid env var %s: %v", config.EnvVarGoplsFlags, err)
	} else {
		goplsArgs = append(goplsArgs, flags...)
	}

	gopls := exec.Command(s.g.goplspath, goplsArgs...)
	gopls.Env = s.g.goplsEnv
	if ev, ok := os.LookupEnv(string(config.EnvVarGoplsGOMAXPROCSMinusN)); ok {
		v := strings.TrimSpace(ev)
		var gmp int
		if strings.HasSuffix(v, "%") {
			v = strings.TrimSuffix(v, "%")
			p, err := strconv.ParseFloat(v, 10)
			if err != nil {
				return fmt.Errorf("failed to parse percentage from %v value %q: %v", config.EnvVarGoplsGOMAXPROCSMinusN, ev, err)
			}
			gmp = int(math.Floor(float64(runtime.NumCPU()) * (1 - p/100)))
		} else {
			n, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("failed to parse integer from %v value %q: %v", config.EnvVarGoplsGOMAXPROCSMinusN, ev, err)
			}
			gmp = runtime.NumCPU() - n
		}
		if gmp < 0 || gmp > runtime.NumCPU() {
			return fmt.Errorf("%v value %q results in GOMAXPROCS value %v which is invalid", config.EnvVarGoplsGOMAXPROCSMinusN, ev, gmp)
		}
		s.g.Logf("Starting gopls with GOMAXPROCS=%v", gmp)
		gopls.Env = append(gopls.Env, "GOMAXPROCS="+strconv.Itoa(gmp))
	}

	s.Logf("Running gopls: %v", strings.Join(gopls.Args, " "))
	stderr, err := gopls.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe for gopls: %v", err)
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			s.g.Logf("gopls stderr: %v", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			s.g.Logf("Error reading gopls stderr: %s", err)
		}
	}()
	stdout, err := gopls.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe for gopls: %v", err)
	}
	stdin, err := gopls.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe for gopls: %v", err)
	}
	s.goplsStdinPipe = stdin
	if err := gopls.Start(); err != nil {
		return fmt.Errorf("failed to start gopls: %v", err)
	}
	go func() {
		err := gopls.Wait()
		if err == nil {
			err = errors.New("gopls exited unexpectedly (status 0)")
		} else {
			err = fmt.Errorf("got error running gopls: %v", err)
		}
		select {
		case s.errCh <- err:
		default:
		}
	}()

	fakeconn := fakenet.NewConn("stdio", stdout, stdin)
	stream := jsonrpc2.NewHeaderStream(fakeconn)
	conn := jsonrpc2.NewConn(stream)
	server := protocol.ServerDispatcher(conn)
	handler := protocol.ClientHandler(s.g, jsonrpc2.MethodNotFound)
	handler = protocol.Handlers(handler)
	ctxt := protocol.WithClient(context.Background(), s.g)

	go func() {
		conn.Go(ctxt, handler)
		<-conn.Done()
		s.g.Logf("fakeconn exited with %s", conn.Err())
	}()

	if _, err := server.Initialize(context.Background(), initParams); err != nil {
		return fmt.Errorf("failed to initialise gopls: %v", err)
	}

	if err := server.Initialized(context.Background(), &protocol.InitializedParams{}); err != nil {
		return fmt.Errorf("failed to call gopls.Initialized: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.u.Store(server)
	s.goplsStdinPipe = stdin

	return nil
}

func (s *goplsServer) stop() error {
	close(s.quit)

	s.mu.Lock()
	defer s.mu.Unlock()
	// We "kill" gopls by closing its stdin. Standard practice for processes
	// that communicate over stdin/stdout is to exit cleanly when stdin is
	// closed.
	if err := s.goplsStdinPipe.Close(); err != nil {
		return err
	}
	return nil
}

func (g *govimplugin) startGopls() error {
	initParams := &protocol.ParamInitialize{}
	initParams.RootURI = protocol.DocumentURI(span.URIFromPath(g.vimstate.workingDirectory))
	initParams.Capabilities.TextDocument.Hover = protocol.HoverClientCapabilities{
		ContentFormat: []protocol.MarkupKind{protocol.PlainText},
	}
	initParams.Capabilities.Workspace.Configuration = true
	// TODO: actually handle these registrations dynamically, if we ever want to
	// target language servers other than gopls.
	initParams.Capabilities.Workspace.DidChangeConfiguration.DynamicRegistration = true
	initParams.Capabilities.Workspace.DidChangeWatchedFiles.DynamicRegistration = true

	initParams.Capabilities.Window.WorkDoneProgress = true

	// Session-level config should be able to be set post initialize, but that
	// is not currently supported by gopls. So for now a restart is required
	// in order to change symbol matcher/style config
	//
	// TODO: clarify whether this method is in fact running as part of the vimstate
	// "thread" and hence whether this lock is required
	g.vimstate.configLock.Lock()
	conf := g.vimstate.config
	defer g.vimstate.configLock.Unlock()
	goplsConfig := make(map[string]interface{})
	if conf.SymbolMatcher != nil {
		goplsConfig[goplsSymbolMatcher] = *conf.SymbolMatcher
	}
	if conf.SymbolStyle != nil {
		goplsConfig[goplsSymbolStyle] = *conf.SymbolStyle
	}

	// TODO: This option was introduced as a way to opt-out from the changes introduced in CL 268597.
	// According to CL 274532 (that added this opt-out), it is intended to be removed - "Ideally
	// we'll be able to remove them in a few months after things stabilize.". We need to handle that
	// case before it is removed.
	goplsConfig["allowModfileModifications"] = true

	initParams.InitializationOptions = goplsConfig

	g.server = &goplsServer{
		g:     g,
		quit:  make(chan struct{}),
		errCh: make(chan error, 1),
	}

	// Initially, gopls must successfully start once.
	if err := g.server.start(initParams); err != nil {
		return err
	}

	gomodpath, err := goModPath(g.vimstate.workingDirectory)
	if err != nil {
		return fmt.Errorf("failed to derive go.mod path: %v", err)
	}

	if gomodpath != "" {
		// i.e. we are in a module
		mw, err := newModWatcher(g, gomodpath)
		if err != nil {
			return fmt.Errorf("failed to create modWatcher for %v: %v", gomodpath, err)
		}
		g.modWatcher = mw
	}

	// Now restart gopls if it crashes.
	go func() {
		for {
			select {
			case <-g.server.quit:
				return
			case err := <-g.server.errCh:
				g.Logf("gopls exited unexpectedly: %s", err)
				params := &protocol.ShowMessageParams{
					Type:    protocol.Error,
					Message: fmt.Sprintf("gopls exited unexpectedly: %s", err),
				}
				g.ShowMessage(context.Background(), params)
				for {
					g.Logf("Restarting gopls...")
					err := g.server.start(initParams)
					if err == nil {
						break
					}
					g.Logf("Failed to restart gopls: %s", err)
					t := time.NewTimer(5 * time.Second)
					defer t.Stop()
					select {
					case <-t.C:
					case <-g.server.quit:
						return
					}
				}
			}
		}
	}()
	return nil
}
