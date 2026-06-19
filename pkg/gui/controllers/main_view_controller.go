package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainViewController struct {
	baseController
	c *ControllerCommon

	context      *context.MainContext
	otherContext *context.MainContext
}

var _ types.IController = &MainViewController{}

func NewMainViewController(
	c *ControllerCommon,
	context *context.MainContext,
	otherContext *context.MainContext,
) *MainViewController {
	return &MainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
		otherContext:   otherContext,
	}
}

func (self *MainViewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:            opts.GetKeys(opts.Config.Universal.TogglePanel),
			Handler:         self.togglePanel,
			Description:     self.c.Tr.ToggleStagingView,
			Tooltip:         self.c.Tr.ToggleStagingViewTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Return),
			Handler:         self.escape,
			Description:     self.c.Tr.ExitFocusedMainView,
			DisplayOnScreen: true,
		},
		{
			// overriding this because we want to read all of the task's output before we start searching
			Keys:        opts.GetKeys(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.StartSearch,
			Tag:         "navigation",
		},
		{
			// These are the same keys we use for jumping to the next/previous
			// search match. While searching, gocui intercepts them for that
			// purpose, so this handler only ever runs when we're not searching.
			Keys:        opts.GetKeys(opts.Config.Universal.NextMatch),
			Handler:     self.handleNextChange,
			Description: self.c.Tr.NextChange,
			Tag:         "navigation",
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.PrevMatch),
			Handler:     self.handlePreviousChange,
			Description: self.c.Tr.PrevChange,
			Tag:         "navigation",
		},
	}
}

func (self *MainViewController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInAlreadyFocusedView,
			FocusedView: self.context.GetViewName(),
		},
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInOtherViewOfMainViewPair,
			FocusedView: self.otherContext.GetViewName(),
		},
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
}

func (self *MainViewController) togglePanel() error {
	if self.otherContext.GetView().Visible {
		self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	}

	return nil
}

func (self *MainViewController) escape() error {
	self.c.Context().Pop()
	return nil
}

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	if self.c.UserConfig().Gui.DisableClickToStageLines {
		return nil
	}

	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(self.context.GetViewName(), opts.Y)
	}
	return nil
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Push(self.context, types.OnFocusOpts{
		ClickedWindowName:  self.context.GetWindowName(),
		ClickedViewLineIdx: opts.Y,
	})

	return nil
}

func (self *MainViewController) openSearch() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				return self.c.Helpers().Search.OpenSearchPrompt(self.context)
			})
		})
	}

	return nil
}

func (self *MainViewController) handleNextChange() error {
	return self.scrollToChange(true)
}

func (self *MainViewController) handlePreviousChange() error {
	return self.scrollToChange(false)
}

// Scrolls the diff in the main view so that the next (or previous) changed
// region is at the top of the viewport.
func (self *MainViewController) scrollToChange(forward bool) error {
	view := self.context.GetView()
	lines := view.ViewBufferLines()
	coloredBackground := view.ViewLinesHaveColoredBackground()
	if len(coloredBackground) != len(lines) {
		return nil
	}

	// A line is part of a change either because it carries a +/- diff marker
	// (the default, non-pager diff) or because a pager such as delta rendered it
	// with a colored background instead of a marker.
	isChange := func(i int) bool {
		return coloredBackground[i] || isDiffChangeMarker(lines[i])
	}

	idx, found := nextChangeRegionIdx(len(lines), view.OriginY(), forward, isChange)
	if !found {
		return nil
	}

	if !view.CanScrollPastBottom {
		maxOriginY := max(len(lines)-view.InnerHeight(), 0)
		idx = min(idx, maxOriginY)
	}
	view.SetOriginY(idx)

	return nil
}

// nextChangeRegionIdx returns the index of the line that starts the next (or
// previous) region of changed lines relative to fromIdx (exclusive), and
// whether one was found. Consecutive changed lines form a single region, so a
// region starts at a changed line whose predecessor is not a change.
func nextChangeRegionIdx(numLines int, fromIdx int, forward bool, isChange func(int) bool) (int, bool) {
	isRegionStart := func(i int) bool {
		return isChange(i) && (i == 0 || !isChange(i-1))
	}

	if forward {
		for i := fromIdx + 1; i < numLines; i++ {
			if isRegionStart(i) {
				return i, true
			}
		}
	} else {
		for i := min(fromIdx, numLines) - 1; i >= 0; i-- {
			if isRegionStart(i) {
				return i, true
			}
		}
	}

	return 0, false
}

// isDiffChangeMarker reports whether a rendered diff line is an added or removed
// line, identified by its leading +/- marker. The +++/--- file-header lines are
// not changes.
func isDiffChangeMarker(line string) bool {
	if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
		return false
	}

	return strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")
}
