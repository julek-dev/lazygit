package ui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableClickToStageLinesInADiffStat = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With disableClickToStageLines set, clicking a file in a focused diff's diffstat scrolls to the file without moving the selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       30,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.DisableClickToStageLines = true
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 600)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%03d", i+1)
		}
		shell.CreateFileAndAdd("aaa.txt", strings.Join(lines, "\n")+"\n")
		shell.CreateFileAndAdd("dir/bbb.txt", "one\n")
		shell.CreateFileAndAdd("zzz.txt", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("one")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+line001")).
			TopLines(
				Contains("commit"),
				Contains("Author:"),
				Contains("Date:"),
				Equals(""),
				Contains("one"),
				Equals("---"),
				Contains("aaa.txt"),
				Contains("dir/bbb.txt"),
				Contains("zzz.txt"),
			).
			Click(2, 8).
			TopVisibleLine(Contains("diff --git a/zzz.txt b/zzz.txt")).
			SelectedLines(Contains("+line001"))
	},
})
