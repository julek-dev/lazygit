package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var IgnoreWhitespaceSubtitle = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The subtitle saying that whitespace is ignored is on every diff that ignores it, and on nothing else",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.IgnoreWhitespaceInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused()

		t.Views().Main().
			Content(Contains("No changed files")).
			/* EXPECTED:
			Subtitle(Equals(""))
			ACTUAL: */
			Subtitle(Equals("(ignoring whitespace)"))

		t.Views().ReflogCommits().
			Focus().
			Lines(
				Contains("commit: second commit").IsSelected(),
				Contains("commit (initial): first commit"),
			)

		t.Views().Main().
			Title(Equals("Reflog Entry")).
			/* EXPECTED:
			Subtitle(Equals("(ignoring whitespace)"))
			ACTUAL: */
			Subtitle(Equals(""))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressPrimaryAction()

		t.Views().Secondary().
			Title(Equals("Custom patch")).
			/* EXPECTED:
			Subtitle(Equals("(ignoring whitespace)"))
			ACTUAL: */
			Subtitle(Equals(""))
	},
})
