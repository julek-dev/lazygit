package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickToSelectLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click on a line in a commit file diff to enter the patch building view at that line, and click inside the patch building view to move the line selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
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

		// clicking a line in the read-only diff enters the patch building view at
		// that line
		t.Views().Main().
			IsFocused().
			Content(Contains("+three")).
			Click(1, 8) // '+four'

		t.Views().PatchBuilding().
			IsFocused().
			SelectedLines(Contains("+four")).
			// clicking another line moves the selection there
			Click(1, 7).
			SelectedLines(Contains("+three"))
	},
})
