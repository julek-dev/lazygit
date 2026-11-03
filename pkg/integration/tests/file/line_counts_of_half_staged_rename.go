package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LineCountsOfHalfStagedRename = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the line counts of a file that is moved with only the new file staged, which git status lists as a deletion and an addition",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowNumstatInFilesView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("add original")

		shell.CreateFileAndAdd("renamed", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\neleven\n")
		shell.DeleteFile("original")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("   D original -10"),
				Equals("  A  renamed +11"),
			)
	},
})
