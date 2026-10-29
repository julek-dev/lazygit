package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LineCountsInCommitFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the number of lines that a commit adds to and deletes from each of its files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowNumstatInCommitFilesView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("modified", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("deleted", "one\ntwo\n")
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("modified", "one\nchanged\nthree\nfour\n")
		shell.DeleteFileAndAdd("deleted")
		shell.CreateDir("dir")
		shell.RunCommand([]string{"git", "mv", "original", "dir/renamed"})
		shell.UpdateFileAndAdd("dir/renamed", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\neleven\n")
		shell.CreateFileAndAdd("dir/added", "alpha\nbeta\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  D deleted -2"),
				Equals("  ▼ dir"),
				Equals("    A added +2"),
				Equals("    R original → renamed +1"),
				Equals("  M modified +2 -1"),
			)
	},
})
