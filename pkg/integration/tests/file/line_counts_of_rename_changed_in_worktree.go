package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LineCountsOfRenameChangedInWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the line counts of a staged rename whose file the worktree has since changed beyond recognition",
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
				Equals("RM original → renamed +3 -5"),
			)
	},
})
