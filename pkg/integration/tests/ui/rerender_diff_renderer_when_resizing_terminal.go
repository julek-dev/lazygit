package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RerenderDiffRendererWhenResizingTerminal = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resizing the terminal re-renders a width-dependent diff renderer at the new width",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        150,
	Height:       100,
	SetupConfig: func(appConfig *config.AppConfig) {
		appConfig.GetUserConfig().Gui.UseHunkModeInDiffView = false
		appConfig.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{
				// Announcing the metadata protocol keeps the renderer's output on
				// screen when the main view is focused, rather than git's own diff.
				Command: `printf '\033]1717;1\007COLUMNWIDTH={{columnWidth}}\n'; cat`,
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")
		shell.UpdateFile("file1", "first line\nsecond line\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1"))

		t.Views().Main().Content(Contains("COLUMNWIDTH=43"))

		// Widening the terminal widens the main view, so the renderer must be re-run
		// so that its columns use the new width.
		t.ResizeTerminal(250, 100)

		t.Views().Main().
			Content(Contains("COLUMNWIDTH=76"))

		// The same goes for the focused main view, and its selection stays on the
		// line it was on.
		t.Views().Files().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+second line")).
			SelectPreviousItem().
			SelectedLines(Contains(" first line"))

		t.ResizeTerminal(200, 100)

		t.Views().Main().
			Content(Contains("COLUMNWIDTH=59")).
			SelectedLines(Contains(" first line"))
	},
})
