package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickToStageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click on a line in a diff to enter the staging view at that line, and click inside the staging view to move the line selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\nfive\nsix\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1")).
			// focus the read-only main diff view
			Press(keys.Universal.FocusMainView)

		// clicking a line in the read-only diff enters the staging view at that line
		t.Views().Main().
			IsFocused().
			Content(Contains("+three")).
			Click(1, 8) // '+four'

		t.Views().Staging().
			IsFocused().
			SelectedLines(Contains("+four")).
			// clicking another line moves the selection there
			Click(1, 7).
			SelectedLines(Contains("+three")).
			// stage two lines so that the staged (secondary) view has content
			PressPrimaryAction().
			PressPrimaryAction().
			SelectedLines(Contains("+five"))

		// clicking an unfocused staging view focuses it and jumps the selection to
		// the clicked line
		t.Views().StagingSecondary().
			Click(1, 8). // '+four'
			IsFocused().
			SelectedLines(Contains("+four"))

		t.Views().Staging().
			Click(1, 9). // '+six'
			IsFocused().
			SelectedLines(Contains("+six"))
	},
})
