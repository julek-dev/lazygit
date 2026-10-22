package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DirectoryHeadingFollowsATerminalResize = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The rule over the files of an untracked directory is drawn again across the view when the terminal gets narrower",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        150,
	Height:       100,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.SetConfig("status.showUntrackedFiles", "normal")
		shell.CreateFile("dir/file1", "content")
		shell.CreateFile("dir/file2", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? dir/").IsSelected(),
			)

		// A rule drawn across a wider view would wrap onto a second line.
		t.Views().Main().
			ContainsViewLines(
				Contains("── Files in 'dir' ─"),
				Equals(""),
			)

		t.ResizeTerminal(100, 100)

		t.Views().Main().
			ContainsViewLines(
				Contains("── Files in 'dir' ─"),
				Equals(""),
			)
	},
})
