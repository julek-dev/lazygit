package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDiffLinesOfTheWholeFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Showing the whole file gives a diff the context to stage lines from, even with a context size of 0",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffContextSize = 0
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")
		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("not possible with a diff context size of 0")).
			Confirm()

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.ToggleWholeFileInDiffView).
			SelectedLines(
				Contains("-three"),
			).
			PressPrimaryAction()

		t.Views().Secondary().
			ContainsLines(
				Contains(" two"),
				Contains("-three"),
				Contains(" four"),
			).
			Content(DoesNotContain("THREE").DoesNotContain("NINE"))
		t.Views().Main().
			ContainsLines(
				Contains(" two"),
				Contains("+THREE"),
				Contains(" four"),
			)
	},
})
