package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LineCountsOfRenameBelowDefaultThreshold = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the line counts of a rename that only the configured similarity threshold detects",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowNumstatInFilesView = true
		config.GetUserConfig().Git.RenameSimilarityThreshold = 45
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("add original")

		shell.DeleteFileAndAdd("original")
		shell.CreateFileAndAdd("renamed", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				/* EXPECTED:
				Equals("R  original → renamed +5"),
				ACTUAL: */
				Equals("R  original → renamed"),
			)
	},
})
