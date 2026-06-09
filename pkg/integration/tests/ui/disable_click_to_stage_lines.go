package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableClickToStageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With disableClickToStageLines set, clicking a diff neither focuses it nor changes its selection; only keys do",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.DisableClickToStageLines = true
		config.GetUserConfig().Gui.ShowFileTree = false
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\n")
		shell.CreateFile("file2", "new\n")
		shell.CreateFile("file3", "new\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file1"))

		t.Views().Main().
			Click(1, 8).
			SelectionIsHidden()

		// A drag that starts in the diff doesn't act on the views it moves over either.
		t.Views().Main().
			ClickAndHold(1, 8).
			MouseMoveToView(t.Views().Files(), 1, 2).
			MouseRelease()

		t.Views().Files().
			IsFocused().
			SelectedLines(Contains("file1")).
			Press(keys.Universal.FocusMainView)

		// Hunk mode is the default, so the whole block of changes is selected, and
		// stays so whatever is clicked, including by the click that focuses the
		// terminal window.
		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+three"), Contains("+four")).
			FocusInAndClick(1, 8).
			SelectedLines(Contains("+three"), Contains("+four")).
			Click(1, 5).
			SelectedLines(Contains("+three"), Contains("+four")).
			Drag(1, 7, 1, 8).
			SelectedLines(Contains("+three"), Contains("+four")).
			// Alt-clicking a line still opens it in the editor.
			AltClick(1, 8).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/repo/file1:4\n"))
			}).
			SelectedLines(Contains("+three"), Contains("+four")).
			// Keys still select lines, and act on them.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(Contains("+three")).
			Click(1, 8).
			SelectedLines(Contains("+three")).
			SelectNextItem().
			SelectedLines(Contains("+four")).
			PressPrimaryAction()

		// The staged changes beside the diff don't take the focus on a click either.
		t.Views().Secondary().
			Content(Contains("+four")).
			Click(1, 5).
			SelectionIsHidden()

		t.Views().Main().
			IsFocused()

		// A main view showing a message rather than a diff still takes the focus.
		t.Views().Stash().
			Focus()

		t.Views().Main().
			Content(DoesNotContain("diff --git")).
			Click(1, 0).
			IsFocused()
	},
})
