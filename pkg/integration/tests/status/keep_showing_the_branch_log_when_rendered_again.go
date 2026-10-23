package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepShowingTheBranchLogWhenRenderedAgain = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The branch log shown with the log key stays on screen when the status panel renders to the main view again, until another panel is focused",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AllBranchesLogCmds = []string{`echo "view1"`, `echo "view2"`}
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus().
			Press(keys.Status.AllBranchesLogGraph).
			Press(keys.Status.AllBranchesLogGraph)

		t.Views().Main().Content(Contains("view2"))

		// A screen mode change renders the panel to the main view again, and so does
		// closing a popup.
		t.Views().Status().
			Press(keys.Universal.NextScreenMode)

		t.Views().Main().
			/* EXPECTED:
			Content(Contains("view2"))
			ACTUAL: */
			Content(Contains("Copyright"))

		t.Views().Status().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Cancel()

		t.Views().Main().
			/* EXPECTED:
			Content(Contains("view2"))
			ACTUAL: */
			Content(Contains("Copyright"))

		// Coming back from another panel shows the panel's default view again.
		t.Views().Files().
			Focus()

		t.Views().Status().
			Focus()

		t.Views().Main().Content(Contains("Copyright"))
	},
})
