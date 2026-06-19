package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NavigateToChangeInMainViewWithPager = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "n/N jump between changes when the diff is rendered by a pager (like delta) that conveys changes through background color rather than +/- markers",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// Stand in for a pager like delta: drop the +/- markers and instead mark
		// changed lines with a background color, so the only way to find changes
		// is by that color.
		cfg.GetUserConfig().Git.Pagers = []config.PagingConfig{
			{
				ColorArg: "never",
				Pager:    `awk '/^\+\+\+/{print;next} /^---/{print;next} /^\+/{printf "\033[48;5;22m%s\033[0m\n",substr($0,2);next} /^-/{printf "\033[48;5;52m%s\033[0m\n",substr($0,2);next} {print}'`,
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nline 10\nline 11\nline 12\nline 13\nline 14\nline 15\nline 16\nline 17\nline 18\nline 19\nline 20\n")
		shell.Commit("initial")
		shell.UpdateFile("file1", "CHANGED FIRST\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nline 10\nline 11\nline 12\nline 13\nline 14\nline 15\nline 16\nline 17\nline 18\nline 19\nCHANGED LAST\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1")).
			Press(keys.Universal.FocusMainView)

		// The pager stripped the +/- markers, so the changed lines ("line 1" and
		// "CHANGED FIRST") are only distinguishable by their background color.
		t.Views().Main().
			IsFocused().
			Press(keys.Universal.NextMatch).
			VisibleLinesFromTop(
				Contains("line 1"),
				Contains("CHANGED FIRST"),
			)

		t.Views().Main().
			Press(keys.Universal.NextMatch).
			VisibleLinesFromTop(
				Contains("line 20"),
				Contains("CHANGED LAST"),
			)

		t.Views().Main().
			Press(keys.Universal.PrevMatch).
			VisibleLinesFromTop(
				Contains("line 1"),
				Contains("CHANGED FIRST"),
			)
	},
})
