package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableClickToSelectLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With disableClickToStageLines set, clicking a commit file diff neither enters the patch building view nor moves the line selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.DisableClickToStageLines = true
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1")).
			// focus the read-only main diff view
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(Contains("+three"))

		// the patch building view never gets populated (we didn't enter it), and
		// focus remains on the read-only diff
		t.Views().Main().Click(1, 7)
		t.Views().PatchBuilding().IsEmpty()
		t.Views().Main().IsFocused()

		// entering the patch building view with the keyboard still works, but
		// clicking inside it must not move the selection
		t.Views().Main().PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			// '+three' is the first changed line and is selected on entry
			SelectedLines(Contains("+three")).
			// clicking '+four' would normally move the selection there; with the
			// option set it's a no-op, so '+three' stays selected
			Click(1, 8).
			SelectedLines(Contains("+three"))
	},
})
