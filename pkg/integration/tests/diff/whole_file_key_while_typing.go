package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var WholeFileKeyWhileTyping = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The key that shows the whole file does nothing while typing in a popup",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("abc")

		t.Views().CommitMessage().
			Press(keys.Universal.ToggleWholeFileInDiffView).
			Content(Equals("abc"))

		t.ExpectPopup().CommitMessagePanel().
			Cancel()

		t.Views().Files().
			IsFocused()
		t.Views().Secondary().
			Subtitle(Equals("")).
			Content(DoesNotContain(" one"))
	},
})
