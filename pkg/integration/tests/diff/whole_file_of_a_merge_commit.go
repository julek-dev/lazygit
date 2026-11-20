package diff

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var WholeFileOfAMergeCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Showing the whole file leaves a merge commit's diff as it is, since its context decides which changes it shows",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.ShowWholeFileInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 40)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		base := strings.Join(lines, "\n") + "\n"
		shell.CreateFileAndAdd("file1", base)
		shell.Commit("base")
		shell.NewBranch("side")
		shell.UpdateFileAndAdd("file1", strings.Replace(base, "line05", "SIDE05", 1))
		shell.Commit("side change")
		shell.Checkout("master")
		shell.UpdateFileAndAdd("file1", strings.Replace(base, "line35", "MAIN35", 1))
		shell.Commit("main change")
		shell.Merge("side")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The merge brought in nothing of its own, so it has no diff to show, however
		// much context there is.
		t.Views().Commits().
			Focus().
			Lines(
				Contains("Merge branch 'side'").IsSelected(),
				Contains("side change"),
				Contains("main change"),
				Contains("base"),
			)

		t.Views().Main().
			Subtitle(Equals("")).
			Content(Contains("file1 |").DoesNotContain("SIDE05").DoesNotContain("MAIN35"))

		t.Views().Commits().
			NavigateToLine(Contains("main change"))

		t.Views().Main().
			Subtitle(Equals("(whole file)")).
			Content(Contains("line01").Contains("+MAIN35").Contains("line40"))
	},
})
