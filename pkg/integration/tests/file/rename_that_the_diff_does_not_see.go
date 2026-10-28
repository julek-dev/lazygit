package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenameThatTheDiffDoesNotSee = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show no line counts for a rename that git status lists, but that the diff counting its lines doesn't pair up",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowNumstatInFilesView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("add original")

		shell.RunCommand([]string{"git", "mv", "original", "renamed"})
		shell.UpdateFile("renamed", "alpha\nbeta\ngamma\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				/* EXPECTED:
				Equals("RM original → renamed"),
				ACTUAL: */
				Equals("RM original → renamed +3"),
			)
	},
})
