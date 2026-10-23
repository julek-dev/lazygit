package diff

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StatFollowsATerminalResize = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A commit's diffstat is laid out again to the width the view has after the terminal gets narrower",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        150,
	Height:       100,
	SetupConfig: func(cfg *config.AppConfig) {
		// git's own diff, so that the render runs as a plain command rather
		// than through a renderer.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Type: "rawGit"},
		}
	},
	SetupRepo: func(shell *Shell) {
		// Enough added lines that git scales the graph to the width it has,
		// rather than drawing one mark per line.
		lines := make([]string, 200)
		for i := range lines {
			lines[i] = "line " + strconv.Itoa(i)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("add file1")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("add file1").IsSelected(),
			)

		t.Views().Main().
			ContainsViewLines(
				Contains("file1 | 200"),
				Contains("1 file changed"),
			)

		t.ResizeTerminal(100, 100)

		t.Views().Main().
			ContainsViewLines(
				Contains("file1 | 200"),
				Contains("1 file changed"),
			)
	},
})
