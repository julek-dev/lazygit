package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResizeTerminalWhileStagingALine = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resizing the terminal while a line is being staged from the focused main view leaves the selection on the next change and the keys working",
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
		shell.UpdateFile("file1", "first line\nsecond line\nthird line\n")

		// Git takes a while to report the status of the working tree, so that the
		// terminal is resized while the refresh that staging a line leads to is
		// still running.
		shell.CreateFile(".git/slow-fsmonitor", "#!/bin/sh\nsleep 0.2\nexit 1\n")
		shell.MakeExecutable(".git/slow-fsmonitor")
		shell.SetConfig("core.fsmonitor", ".git/slow-fsmonitor")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file1")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+second line")).
			PressAndResizeTerminal(keys.Universal.Select, 200, 100).
			Content(Contains("COLUMNWIDTH=34")).
			SelectedLines(Contains("+third line")).
			SelectPreviousItem().
			SelectedLines(Contains(" second line"))

		t.Views().Secondary().
			Content(Contains("COLUMNWIDTH=34").Contains("+second line"))
	},
})
