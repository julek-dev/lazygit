package controllers

import (
	"fmt"
	"math"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller lets you change the context size for diffs. The 'context' in 'context size' refers to the conventional meaning of the word 'context' in a diff, as opposed to lazygit's own idea of a 'context'.

type ContextLinesController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &ContextLinesController{}

func NewContextLinesController(
	c *ControllerCommon,
) *ContextLinesController {
	return &ContextLinesController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *ContextLinesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:        opts.GetKeys(opts.Config.Universal.IncreaseContextInDiffView),
			Handler:     self.Increase,
			Description: self.c.Tr.IncreaseContextInDiffView,
			Tooltip:     self.c.Tr.IncreaseContextInDiffViewTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.DecreaseContextInDiffView),
			Handler:     self.Decrease,
			Description: self.c.Tr.DecreaseContextInDiffView,
			Tooltip:     self.c.Tr.DecreaseContextInDiffViewTooltip,
		},
	}

	return bindings
}

func (self *ContextLinesController) Context() types.Context {
	return nil
}

func (self *ContextLinesController) Increase() error {
	step := self.c.UserConfig().Git.DiffContextSizeStep
	if self.c.UserConfig().Git.DiffContextSize <= math.MaxUint64-step {
		self.c.UserConfig().Git.DiffContextSize += step
	} else {
		self.c.UserConfig().Git.DiffContextSize = math.MaxUint64
	}
	return self.applyChange()
}

func (self *ContextLinesController) Decrease() error {
	gitConfig := &self.c.UserConfig().Git
	// Stop at the minimum, and leave a size the config puts below it alone.
	if gitConfig.DiffContextSize > gitConfig.DiffContextSizeMin {
		gitConfig.DiffContextSize -= min(gitConfig.DiffContextSizeStep, gitConfig.DiffContextSize-gitConfig.DiffContextSizeMin)
	}
	return self.applyChange()
}

func (self *ContextLinesController) applyChange() error {
	self.c.Toast(fmt.Sprintf(self.c.Tr.DiffContextSizeChanged, self.c.UserConfig().Git.DiffContextSize))

	currentContext := self.c.Context().CurrentSide()
	// The diff is about to be rendered again with more or less context around
	// each change, which reads as the lines you were looking at moving up or down
	// the view; keep them where they are instead.
	self.c.Helpers().DiffLine.PreserveDiffPositionOnRerender(self.c.Contexts().Normal.GetView())
	self.c.Helpers().DiffLine.PreserveDiffPositionOnRerender(self.c.Contexts().NormalSecondary.GetView())
	currentContext.HandleRenderToMain()
	return nil
}
