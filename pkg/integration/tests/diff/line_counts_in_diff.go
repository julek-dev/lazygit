package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LineCountsInDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the line counts of the files of a commit while diffing against a branch, in both directions",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowNumstatInCommitFilesView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\nthird line\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
				Contains("branch-b"),
			).
			Press(keys.Universal.DiffingMenu)

		t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Contains(`Diff branch-a`)).Confirm()

		t.Views().Branches().
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("update")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M file1 +2"),
			).
			PressEscape()

		t.Views().SubCommits().
			IsFocused().
			PressEscape()

		t.Views().Branches().
			IsFocused().
			Press(keys.Universal.DiffingMenu)

		t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Contains("Reverse diff direction")).Confirm()

		t.Views().Branches().
			IsFocused().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("update")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M file1 -2"),
			)
	},
})
