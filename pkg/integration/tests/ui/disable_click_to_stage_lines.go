package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableClickToStageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With disableClickToStageLines set, clicking a diff neither enters the staging view nor moves the line selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.DisableClickToStageLines = true
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\nfive\nsix\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Layer A: clicking a line in the read-only diff must not enter the staging view.
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1")).
			// focus the read-only main diff view
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(Contains("+three"))

		// the staging view never gets populated (we didn't enter it), and focus
		// remains on the read-only diff
		t.Views().Main().Click(1, 7)
		t.Views().Staging().IsEmpty()
		t.Views().Main().IsFocused()

		// Layer B: clicking inside the staging view must not move the selection.
		t.Views().Files().
			Focus().
			PressEnter()

		t.Views().Staging().
			IsFocused().
			// '+three' is the first changed line and is selected on entry
			SelectedLines(Contains("+three")).
			// clicking '+four' would normally move the selection there; with the
			// option set it's a no-op, so '+three' stays selected
			Click(1, 8).
			SelectedLines(Contains("+three")).
			// stage two lines (the keyboard still works) so that the staged
			// (secondary) view has content
			PressPrimaryAction().
			PressPrimaryAction().
			SelectedLines(Contains("+five"))

		// Layer C: clicking an unfocused staging view still focuses it, but doesn't
		// jump the selection to the clicked line.
		t.Views().StagingSecondary().
			// click '+four'; the selection stays on '+three', the first changed line
			Click(1, 8).
			IsFocused().
			SelectedLines(Contains("+three"))

		t.Views().Staging().
			// click '+six'; the selection stays on '+five', the first changed line
			Click(1, 9).
			IsFocused().
			SelectedLines(Contains("+five"))
	},
})
