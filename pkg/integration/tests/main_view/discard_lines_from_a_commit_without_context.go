package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardLinesFromACommitWithoutContext = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard lines from a commit whose diff is shown without any context",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffContextSize = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\nFOUR\nfive\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+TWO")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Discard lines from commit")).
					Content(Contains("Are you sure you want to discard the selected lines from this commit?")).
					Confirm()
			})

		t.Views().Main().
			Content(DoesNotContain("TWO").Contains("+FOUR"))
		t.Views().Files().IsEmpty()
	},
})
