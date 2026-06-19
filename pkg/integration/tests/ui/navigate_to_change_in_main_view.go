package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NavigateToChangeInMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "In a read-only diff view, n/N jump to the next/previous change when not searching",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nline 10\nline 11\nline 12\nline 13\nline 14\nline 15\nline 16\nline 17\nline 18\nline 19\nline 20\n")
		shell.Commit("initial")
		// Change the first and last lines, producing a diff with two separate hunks.
		shell.UpdateFile("file1", "CHANGED FIRST\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nline 10\nline 11\nline 12\nline 13\nline 14\nline 15\nline 16\nline 17\nline 18\nline 19\nCHANGED LAST\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1")).
			Press(keys.Universal.FocusMainView)

		// The diff starts scrolled to the top, above the first change.
		t.Views().Main().
			IsFocused().
			VisibleLinesFromTop(Contains("diff --git a/file1 b/file1"))

		// Jump to the first change.
		t.Views().Main().
			Press(keys.Universal.NextMatch).
			VisibleLinesFromTop(
				Contains("-line 1"),
				Contains("+CHANGED FIRST"),
			)

		// Jump to the second change.
		t.Views().Main().
			Press(keys.Universal.NextMatch).
			VisibleLinesFromTop(
				Contains("-line 20"),
				Contains("+CHANGED LAST"),
			)

		// There is no further change, so we stay put.
		t.Views().Main().
			Press(keys.Universal.NextMatch).
			VisibleLinesFromTop(
				Contains("-line 20"),
				Contains("+CHANGED LAST"),
			)

		// Jump back to the first change.
		t.Views().Main().
			Press(keys.Universal.PrevMatch).
			VisibleLinesFromTop(
				Contains("-line 1"),
				Contains("+CHANGED FIRST"),
			)

		// There is no previous change, so we stay put.
		t.Views().Main().
			Press(keys.Universal.PrevMatch).
			VisibleLinesFromTop(
				Contains("-line 1"),
				Contains("+CHANGED FIRST"),
			)
	},
})
