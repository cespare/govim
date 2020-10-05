package main

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/protocol"
	"github.com/kr/pretty"
)

type goplsServer struct {
	u atomic.Value // of protocol.Server
	g *govimplugin

	quit  chan struct{} // closed to signal shutdown
	errCh chan error    // gopls exit error

	// mu is held when writing u or when accessing goplsStdinPipe.
	mu             sync.Mutex
	goplsStdinPipe io.WriteCloser
}

var _ protocol.Server = &goplsServer{}

func (s *goplsServer) get() protocol.Server {
	v := s.u.Load()
	if v == nil {
		panic("goplsServer accessed before start")
	}
	return v.(protocol.Server)
}

func (s *goplsServer) Logf(format string, args ...interface{}) {
	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	s.g.Logf("gopls server start =======================\n"+format+"gopls server end =======================\n", args...)
}

func (s *goplsServer) Initialize(ctxt context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	s.Logf("gopls.Initialize() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Initialize(ctxt, params)
	s.Logf("gopls.Initialize() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Initialized(ctxt context.Context, params *protocol.InitializedParams) error {
	s.Logf("gopls.Initialized() call; params:\n%v", pretty.Sprint(params))
	err := s.get().Initialized(ctxt, params)
	s.Logf("gopls.Initialized() return; err: %v", err)
	return err
}

func (s *goplsServer) Shutdown(ctxt context.Context) error {
	s.Logf("gopls.Shutdown() call")
	err := s.get().Shutdown(ctxt)
	s.Logf("gopls.Shutdown() return; err: %v", err)
	return err
}

func (s *goplsServer) Exit(ctxt context.Context) error {
	s.Logf("gopls.Exit() call")
	err := s.get().Exit(ctxt)
	s.Logf("gopls.Exit() return; err: %v", err)
	return err
}

func (s *goplsServer) DidChangeWorkspaceFolders(ctxt context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	s.Logf("gopls.DidChangeWorkspaceFolders() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidChangeWorkspaceFolders(ctxt, params)
	s.Logf("gopls.DidChangeWorkspaceFolders() return; err: %v", err)
	return err
}

func (s *goplsServer) DidChangeConfiguration(ctxt context.Context, params *protocol.DidChangeConfigurationParams) error {
	s.Logf("gopls.DidChangeConfiguration() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidChangeConfiguration(ctxt, params)
	s.Logf("gopls.DidChangeConfiguration() return; err: %v", err)
	return err
}

func (s *goplsServer) DidChangeWatchedFiles(ctxt context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	s.Logf("gopls.DidChangeWatchedFiles() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidChangeWatchedFiles(ctxt, params)
	s.Logf("gopls.DidChangeWatchedFiles() return; err: %v", err)
	return err
}

func (s *goplsServer) Symbol(ctxt context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	s.Logf("gopls.Symbol() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Symbol(ctxt, params)
	s.Logf("gopls.Symbol() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) ExecuteCommand(ctxt context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	s.Logf("gopls.ExecuteCommand() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().ExecuteCommand(ctxt, params)
	s.Logf("gopls.ExecuteCommand() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) DidOpen(ctxt context.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.Logf("gopls.DidOpen() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidOpen(ctxt, params)
	s.Logf("gopls.DidOpen() return; err: %v", err)
	return err
}

func (s *goplsServer) DidChange(ctxt context.Context, params *protocol.DidChangeTextDocumentParams) error {
	s.Logf("gopls.DidChange() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidChange(ctxt, params)
	s.Logf("gopls.DidChange() return; err: %v", err)
	return err
}

func (s *goplsServer) WillSave(ctxt context.Context, params *protocol.WillSaveTextDocumentParams) error {
	s.Logf("gopls.WillSave() call; params:\n%v", pretty.Sprint(params))
	err := s.get().WillSave(ctxt, params)
	s.Logf("gopls.WillSave() return; err: %v", err)
	return err
}

func (s *goplsServer) WillSaveWaitUntil(ctxt context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	s.Logf("gopls.WillSaveWaitUntil() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().WillSaveWaitUntil(ctxt, params)
	s.Logf("gopls.WillSaveWaitUntil() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) DidSave(ctxt context.Context, params *protocol.DidSaveTextDocumentParams) error {
	s.Logf("gopls.DidSave() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidSave(ctxt, params)
	s.Logf("gopls.DidSave() return; err: %v", err)
	return err
}

func (s *goplsServer) DidClose(ctxt context.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.Logf("gopls.DidClose() call; params:\n%v", pretty.Sprint(params))
	err := s.get().DidClose(ctxt, params)
	s.Logf("gopls.DidClose() return; err: %v", err)
	return err
}

func (s *goplsServer) Completion(ctxt context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	s.Logf("gopls.Completion() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Completion(ctxt, params)
	s.Logf("gopls.Completion() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Resolve(ctxt context.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	s.Logf("gopls.Resolve() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Resolve(ctxt, params)
	s.Logf("gopls.Resolve() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Hover(ctxt context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	s.Logf("gopls.Hover() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Hover(ctxt, params)
	s.Logf("gopls.Hover() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) SignatureHelp(ctxt context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	s.Logf("gopls.SignatureHelp() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().SignatureHelp(ctxt, params)
	s.Logf("gopls.SignatureHelp() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Definition(ctxt context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	s.Logf("gopls.Definition() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Definition(ctxt, params)
	s.Logf("gopls.Definition() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) TypeDefinition(ctxt context.Context, params *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	s.Logf("gopls.TypeDefinition() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().TypeDefinition(ctxt, params)
	s.Logf("gopls.TypeDefinition() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Implementation(ctxt context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) {
	s.Logf("gopls.Implementation() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Implementation(ctxt, params)
	s.Logf("gopls.Implementation() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) References(ctxt context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	s.Logf("gopls.References() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().References(ctxt, params)
	s.Logf("gopls.References() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) DocumentHighlight(ctxt context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	s.Logf("gopls.DocumentHighlight() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().DocumentHighlight(ctxt, params)
	s.Logf("gopls.DocumentHighlight() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) DocumentSymbol(ctxt context.Context, params *protocol.DocumentSymbolParams) ([]interface{}, error) {
	s.Logf("gopls.DocumentSymbol() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().DocumentSymbol(ctxt, params)
	s.Logf("gopls.DocumentSymbol() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) CodeAction(ctxt context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	s.Logf("gopls.CodeAction() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().CodeAction(ctxt, params)
	s.Logf("gopls.CodeAction() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) CodeLens(ctxt context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	s.Logf("gopls.CodeLens() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().CodeLens(ctxt, params)
	s.Logf("gopls.CodeLens() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) ResolveCodeLens(ctxt context.Context, params *protocol.CodeLens) (*protocol.CodeLens, error) {
	s.Logf("gopls.ResolveCodeLens() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().ResolveCodeLens(ctxt, params)
	s.Logf("gopls.ResolveCodeLens() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) DocumentLink(ctxt context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	s.Logf("gopls.DocumentLink() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().DocumentLink(ctxt, params)
	s.Logf("gopls.DocumentLink() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) ResolveDocumentLink(ctxt context.Context, params *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	s.Logf("gopls.ResolveDocumentLink() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().ResolveDocumentLink(ctxt, params)
	s.Logf("gopls.ResolveDocumentLink() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) DocumentColor(ctxt context.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	s.Logf("gopls.DocumentColor() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().DocumentColor(ctxt, params)
	s.Logf("gopls.DocumentColor() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) ColorPresentation(ctxt context.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	s.Logf("gopls.ColorPresentation() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().ColorPresentation(ctxt, params)
	s.Logf("gopls.ColorPresentation() return; err: %v; res:\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Formatting(ctxt context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	s.Logf("gopls.Formatting() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Formatting(ctxt, params)
	s.Logf("gopls.Formatting() return; err: %v; res:\n%v\n", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) RangeFormatting(ctxt context.Context, params *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	s.Logf("gopls.RangeFormatting() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().RangeFormatting(ctxt, params)
	s.Logf("gopls.RangeFormatting() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) OnTypeFormatting(ctxt context.Context, params *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	s.Logf("gopls.OnTypeFormatting() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().OnTypeFormatting(ctxt, params)
	s.Logf("gopls.OnTypeFormatting() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Rename(ctxt context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	s.Logf("gopls.Rename() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Rename(ctxt, params)
	s.Logf("gopls.Rename() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) FoldingRange(ctxt context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	s.Logf("gopls.FoldingRange() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().FoldingRange(ctxt, params)
	s.Logf("gopls.FoldingRange() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) Declaration(ctxt context.Context, params *protocol.DeclarationParams) (protocol.Declaration, error) {
	s.Logf("gopls.Declaration() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Declaration(ctxt, params)
	s.Logf("gopls.Declaration() return; err: %v; res\n%v%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) LogTrace(ctxt context.Context, params *protocol.LogTraceParams) error {
	s.Logf("gopls.LogTrace() call; params:\n%v", pretty.Sprint(params))
	err := s.get().LogTrace(ctxt, params)
	s.Logf("gopls.LogTrace() return; err: %v", err)
	return err
}

func (s *goplsServer) PrepareRename(ctxt context.Context, params *protocol.PrepareRenameParams) (*protocol.Range, error) {
	s.Logf("gopls.PrepareRename() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().PrepareRename(ctxt, params)
	s.Logf("gopls.PrepareRename() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) SetTrace(ctxt context.Context, params *protocol.SetTraceParams) error {
	s.Logf("gopls.SetTrace() call; params:\n%v", pretty.Sprint(params))
	err := s.get().SetTrace(ctxt, params)
	s.Logf("gopls.SetTrace() return; err: %v", err)
	return err
}

func (s *goplsServer) SelectionRange(ctxt context.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	s.Logf("gopls.SelectionRange() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().SelectionRange(ctxt, params)
	s.Logf("gopls.SelectionRange() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) NonstandardRequest(ctxt context.Context, method string, params interface{}) (interface{}, error) {
	s.Logf("gopls.NonstandardRequest() call; method: %v, params:\n%v", method, pretty.Sprint(params))
	res, err := s.get().NonstandardRequest(ctxt, method, params)
	s.Logf("gopls.NonstandardRequest() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) IncomingCalls(ctxt context.Context, params *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	s.Logf("gopls.IncomingCalls() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().IncomingCalls(ctxt, params)
	s.Logf("gopls.IncomingCalls() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) OutgoingCalls(ctxt context.Context, params *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	s.Logf("gopls.OutgoingCalls() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().OutgoingCalls(ctxt, params)
	s.Logf("gopls.OutgoingCalls() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) PrepareCallHierarchy(ctxt context.Context, params *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	s.Logf("gopls.PrepareCallHierarchy() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().PrepareCallHierarchy(ctxt, params)
	s.Logf("gopls.PrepareCallHierarchy() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) SemanticTokensFull(ctxt context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	s.Logf("gopls.SemanticTokensFull() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().SemanticTokensFull(ctxt, params)
	s.Logf("gopls.SemanticTokensFull() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) SemanticTokensFullDelta(ctxt context.Context, params *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	s.Logf("gopls.SemanticTokensFullDelta() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().SemanticTokensFullDelta(ctxt, params)
	s.Logf("gopls.SemanticTokensFullDelta() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) SemanticTokensRange(ctxt context.Context, params *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	s.Logf("gopls.SemanticTokensRange() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().SemanticTokensRange(ctxt, params)
	s.Logf("gopls.SemanticTokensRange() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) SemanticTokensRefresh(ctxt context.Context) error {
	s.Logf("gopls.SemanticTokensRefresh() call\n")
	err := s.get().SemanticTokensRefresh(ctxt)
	s.Logf("gopls.SemanticTokensRefresh() return; err: %v", err)
	return err
}

func (s *goplsServer) WorkDoneProgressCancel(ctxt context.Context, params *protocol.WorkDoneProgressCancelParams) error {
	s.Logf("gopls.WorkDoneProgressCancel() call; params:\n%v", pretty.Sprint(params))
	err := s.get().WorkDoneProgressCancel(ctxt, params)
	s.Logf("gopls.WorkDoneProgressCancel() return; err: %v\n", err)
	return err
}

func (s *goplsServer) Moniker(ctxt context.Context, params *protocol.MonikerParams) ([]protocol.Moniker /*Moniker[] | null*/, error) {
	s.Logf("gopls.Moniker() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().Moniker(ctxt, params)
	s.Logf("gopls.Moniker() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}

func (s *goplsServer) ResolveCodeAction(ctxt context.Context, params *protocol.CodeAction) (*protocol.CodeAction, error) {
	s.Logf("gopls.ResolveCodeAction() call; params:\n%v", pretty.Sprint(params))
	res, err := s.get().ResolveCodeAction(ctxt, params)
	s.Logf("gopls.ResolveCodeAction() return; err: %v; res\n%v", err, pretty.Sprint(res))
	return res, err
}
