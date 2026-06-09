package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffContextSizeMin = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The diff context size cannot be decreased below the configured minimum",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.DiffContextSizeMin = 2
	},
	SetupRepo: func(shell *Shell) {
		// a couple of lines of context so the staging view has something to show
		shell.CreateFileAndAdd("file1", "1a\n2a\n3a\n4a\n5a\n6a\n7a\n8a\n9a\n10a\n11a\n12a\n13a\n14a\n15a")
		shell.Commit("one")

		shell.UpdateFile("file1", "1a\n2a\n3b\n4a\n5a\n6a\n7a\n8a\n9a\n10a\n11a\n12a\n13b\n14a\n15a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file1")).
			PressEnter()

		// The default context size is 3 and the configured minimum is 2.
		t.Views().Staging().
			IsFocused().
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			// A further decrease is clamped to the minimum instead of going lower.
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			// Increasing again confirms the size really sat at the minimum.
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 3"))
			})
	},
})
