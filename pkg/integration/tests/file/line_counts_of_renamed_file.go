package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LineCountsOfRenamedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the line counts of a renamed file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowNumstatInFilesView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("add original")

		shell.RunCommand([]string{"git", "mv", "original", "renamed"})
		shell.UpdateFileAndAdd("renamed", "one\ntwo\nthree\nfour\nfive\nsix\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				/* EXPECTED:
				Equals("R  original → renamed +1"),
				ACTUAL: */
				Equals("R  original → renamed"),
			)
	},
})
