package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDiffLinesWithoutContext = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage and discard lines of a working tree's diff that is shown without any context",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffContextSize = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")
		shell.UpdateFile("file1", "one\ntwo\nADDED\nthree\nfour\nfive\nsix\nseven\neight\nOTHER\nnine\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(DoesNotContain(" three")).
			SelectedLines(
				Contains("+ADDED"),
			).
			PressPrimaryAction()

		t.Views().Secondary().
			ContainsLines(
				Contains("+ADDED"),
			).
			Content(DoesNotContain("OTHER"))

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+OTHER"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Discard change")).
					Content(Contains("Are you sure you want to discard this change")).
					Confirm()
			})

		t.Views().Files().Lines(
			Contains("M  file1"),
		)
	},
})
