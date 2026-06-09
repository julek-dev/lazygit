package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffContextChangeBelowMin = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Change the number of diff context lines when starting below the configured minimum",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.DiffContextSize = 0
		config.GetUserConfig().Git.DiffContextSizeMin = 2
	},
	SetupRepo: func(shell *Shell) {
		// need to be working with a few lines so that git perceives it as two separate hunks
		shell.CreateFileAndAdd("file1", "1a\n2a\n3a\n4a\n5a\n6a\n7a\n8a\n9a\n10a\n11a\n12a\n13a\n14a\n15a")
		shell.Commit("one")

		shell.UpdateFile("file1", "1a\n2a\n3b\n4a\n5a\n6a\n7a\n8a\n9a\n10a\n11a\n12a\n13b\n14a\n15a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			SelectedLines(
				Contains(`-3a`),
				Contains(`+3b`),
			).
			// increasing is not affected by the minimum: the size may stay below it
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 1"))
			}).
			SelectedLines(
				Contains(`-3a`),
				Contains(`+3b`),
			).
			// decreasing from below the minimum snaps up to the minimum
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			SelectedLines(
				Contains(`-3a`),
				Contains(`+3b`),
			).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			SelectedLines(
				Contains(`-3a`),
				Contains(`+3b`),
			)
	},
})
