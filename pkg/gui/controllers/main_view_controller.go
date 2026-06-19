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
// region is at the top of the viewport. A changed region begins at each hunk
// header, so we navigate between those.
func (self *MainViewController) scrollToChange(forward bool) error {
	view := self.context.GetView()
	lines := view.ViewBufferLines()

	idx, found := adjacentHunkHeaderIdx(lines, view.OriginY(), forward)
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

// Returns the index of the closest hunk header before or after fromIdx
// (exclusive), depending on the direction, and whether one was found.
func adjacentHunkHeaderIdx(lines []string, fromIdx int, forward bool) (int, bool) {
	if forward {
		for i := fromIdx + 1; i < len(lines); i++ {
			if isHunkHeader(lines[i]) {
				return i, true
			}
		}
	} else {
		for i := min(fromIdx, len(lines)) - 1; i >= 0; i-- {
			if isHunkHeader(lines[i]) {
				return i, true
			}
		}
	}

	return 0, false
}

func isHunkHeader(line string) bool {
	return strings.HasPrefix(line, "@@")
}
