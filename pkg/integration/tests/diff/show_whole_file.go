package diff

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowWholeFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle showing the whole file in the diff, and leave it with the context size keys",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 20)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		lines[9] = "LINE10"
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Main().
			Subtitle(Equals("")).
			Content(DoesNotContain("line01").DoesNotContain("line20"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ToggleWholeFileInDiffView)

		t.Views().Main().
			Subtitle(Equals("(whole file)")).
			Content(Contains("line01").Contains("line20"))

		t.Views().Files().
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Already showing the whole file"))
			})

		t.Views().Main().
			Subtitle(Equals("(whole file)")).
			Content(Contains("line01").Contains("line20"))

		// Less than the whole file is the usual context size, not one line less of it.
		t.Views().Files().
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 3"))
			})

		t.Views().Main().
			Subtitle(Equals("")).
			ContainsLines(
				Contains(" line07"),
				Contains(" line08"),
				Contains(" line09"),
				Contains("-line10"),
				Contains("+LINE10"),
				Contains(" line11"),
				Contains(" line12"),
				Contains(" line13"),
			).
			Content(DoesNotContain("line01").DoesNotContain("line20"))

		t.Views().Files().
			Press(keys.Universal.ToggleWholeFileInDiffView)

		t.Views().Main().
			Subtitle(Equals("(whole file)")).
			Content(Contains("line01").Contains("line20"))

		t.Views().Files().
			Press(keys.Universal.ToggleWholeFileInDiffView)

		t.Views().Main().
			Subtitle(Equals("")).
			Content(DoesNotContain("line01").DoesNotContain("line20"))
	},
})
