package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NextFileSelectsAChangeNotYetRead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Going to the next file in hunk mode selects its first change, even when so much context comes before it that it hasn't been read yet",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
		cfg.GetUserConfig().Git.DiffContextSize = 1000
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 600)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%03d", i+1)
		}
		long := strings.Join(lines, "\n") + "\n"
		shell.CreateFileAndAdd("aaa.txt", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("bbb.txt", long)
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("aaa.txt", "one\nTWO\nthree\n")
		shell.UpdateFileAndAdd("bbb.txt", strings.Replace(long, "line550", "LINE550", 1))
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			Press(keys.Main.NextFile).
			SelectedLines(
				/* EXPECTED:
				Contains("-line550"),
				Contains("+LINE550"),
				ACTUAL: */
				Contains("-two"),
				Contains("+TWO"),
			)
	},
})
