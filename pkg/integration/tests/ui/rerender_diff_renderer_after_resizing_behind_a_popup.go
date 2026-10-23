package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RerenderDiffRendererAfterResizingBehindAPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resizing the terminal while a popup is open re-renders a width-dependent diff renderer once the popup is gone",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        150,
	Height:       100,
	SetupConfig: func(appConfig *config.AppConfig) {
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
		t.Views().Main().Content(Contains("COLUMNWIDTH=43"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		// Nothing is rendered into the main view while a popup has the focus.
		t.ResizeTerminal(250, 100)

		t.Views().Main().Content(Contains("COLUMNWIDTH=43"))

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Cancel()

		t.Views().Main().Content(Contains("COLUMNWIDTH=76"))
	},
})
