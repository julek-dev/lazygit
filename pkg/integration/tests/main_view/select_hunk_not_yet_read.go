package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectHunkNotYetRead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggling hunk selection on context selects the change block below it, even when so much context comes before it that it hasn't been read yet",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        160,
	Height:       50,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffContextSize = 3000
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 3000)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%04d", i+1)
		}
		long := strings.Join(lines, "\n") + "\n"
		shell.CreateFileAndAdd("file1", long)
		shell.Commit("first commit")
		changed := strings.Replace(long, "line0005", "LINE0005", 1)
		changed = strings.Replace(changed, "line2800", "LINE2800", 1)
		shell.UpdateFileAndAdd("file1", changed)
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-line0005"),
			).
			NavigateToLine(Contains(" line0100")).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				/* EXPECTED:
				Contains("-line2800"),
				Contains("+LINE2800"),
				ACTUAL: */
				Contains("-line0005"),
				Contains("+LINE0005"),
			)
	},
})
