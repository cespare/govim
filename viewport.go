package govim

import (
	"encoding/json"
	"fmt"
)

const (
	sysFuncOnViewportChange = sysFuncPref + "OnViewportChange"
)

// SubOnViewportChange creates a subscription to the OnViewportChange event
// exposed by Govim
func (g *Govim) SubOnViewportChange(f func(Viewport)) *OnViewportChangeSub {
	res := &OnViewportChangeSub{f: f}
	g.onViewportChangeSubsLock.Lock()
	g.onViewportChangeSubs = append(g.onViewportChangeSubs, res)
	g.onViewportChangeSubsLock.Unlock()
	return res
}

// UnsubOnViewportChange removes a subscription to the OnViewportChange event.
// It panics if sub is not an active subscription.
func (g *Govim) UnsubOnViewportChange(sub *OnViewportChangeSub) {
	g.onViewportChangeSubsLock.Lock()
	defer g.onViewportChangeSubsLock.Unlock()
	for i, s := range g.onViewportChangeSubs {
		if sub == s {
			g.onViewportChangeSubs = append(g.onViewportChangeSubs[:i], g.onViewportChangeSubs[i+1:]...)
			return
		}
	}
	panic(fmt.Errorf("did not find subscription"))
}

type OnViewportChangeSub struct {
	f func(Viewport)
}

func (g *Govim) onViewportChange(args ...json.RawMessage) (interface{}, error) {
	var r Viewport
	g.decodeJSON(args[0], &r)
	g.viewportLock.Lock()
	g.currViewport = r
	g.viewportLock.Unlock()

	var subs []*OnViewportChangeSub
	r = r.dup()
	g.onViewportChangeSubsLock.Lock()
	subs = append(subs, g.onViewportChangeSubs...)
	g.onViewportChangeSubsLock.Unlock()
	for _, s := range subs {
		s.f(r)
	}
	return nil, nil
}

type Viewport struct {
	TabNr   int
	Windows []WinInfo
}

type WinInfo struct {
	WinNr    int
	BotLine  int
	Height   int
	BufNr    int
	WinBar   int
	Width    int
	TabNr    int
	QuickFix bool
	TopLine  int
	LocList  bool
	WinCol   int
	WinRow   int
	WinID    int
	Terminal bool
}

// Viewport returns the active Vim viewport
func (g *Govim) Viewport() Viewport {
	var res Viewport
	g.viewportLock.Lock()
	res = g.currViewport.dup()
	g.viewportLock.Unlock()
	return res
}

func (v Viewport) dup() Viewport {
	v.Windows = append([]WinInfo{}, v.Windows...)
	return v
}

func (wi *WinInfo) UnmarshalJSON(b []byte) error {
	var w struct {
		WinNr    int `json:"winnr"`
		BotLine  int `json:"botline"`
		Height   int `json:"height"`
		BufNr    int `json:"bufnr"`
		WinBar   int `json:"winbar"`
		Width    int `json:"width"`
		TabNr    int `json:"tabnr"`
		QuickFix int `json:"quickfix"`
		TopLine  int `json:"topline"`
		LocList  int `json:"loclist"`
		WinCol   int `json:"wincol"`
		WinRow   int `json:"winrow"`
		WinID    int `json:"winid"`
		Terminal int `json:"terminal"`
	}

	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	wi.WinNr = w.WinNr
	wi.BotLine = w.BotLine
	wi.Height = w.Height
	wi.BufNr = w.BufNr
	wi.WinBar = w.WinBar
	wi.Width = w.Width
	wi.TabNr = w.TabNr
	wi.QuickFix = w.QuickFix == 1
	wi.TopLine = w.TopLine
	wi.LocList = w.LocList == 1
	wi.WinCol = w.WinCol
	wi.WinRow = w.WinRow
	wi.WinID = w.WinID
	wi.Terminal = w.Terminal == 1

	return nil
}
