package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableClickToStageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With disableClickToStageLines set, clicking a diff focuses it without moving the line selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.DisableClickToStageLines = true
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file1"))

		// Clicking '+four' still focuses the diff, but the selection starts on the
		// first change, as it does when focusing the diff with the keyboard.
		t.Views().Main().
			Click(1, 8).
			IsFocused().
			SelectedLines(Contains("+three")).
			// Clicking or dragging inside the focused diff leaves the selection alone.
			Click(1, 8).
			SelectedLines(Contains("+three")).
			Drag(1, 7, 1, 8).
			SelectedLines(Contains("+three")).
			// Alt-clicking a line still opens it in the editor.
			AltClick(1, 8).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/repo/file1:4\n"))
			}).
			SelectedLines(Contains("+three")).
			// The keyboard still moves the selection, and acts on it.
			SelectNextItem().
			SelectedLines(Contains("+four")).
			PressPrimaryAction()

		// Clicking the context line ' one' of the staged changes focuses them, again
		// with the selection on the first change.
		t.Views().Secondary().
			Click(1, 5).
			IsFocused().
			SelectedLines(Contains("+four"))
	},
})
