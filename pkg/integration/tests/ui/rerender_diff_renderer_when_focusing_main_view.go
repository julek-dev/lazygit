package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Diff renderers can produce output that depends on the width of the view, e.g.
// a side-by-side diff whose column width is derived from {{columnWidth}}. When
// the main view is focused in half or full screen mode it widens (the side
// panels are hidden), so the renderer must be re-run with the new width;
// otherwise it keeps showing the narrow output it produced while the side panels
// were taking up space.

var RerenderDiffRendererWhenFocusingMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view in full screen mode re-renders a width-dependent diff renderer at the new width",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(appConfig *config.AppConfig) {
		appConfig.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{
				// Announcing the metadata protocol keeps the renderer's output on
				// screen when the main view is focused, rather than git's own diff.
				Command: `printf '\033]1717;1\007COLUMNWIDTH={{columnWidth}}\n'; cat`,
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")
		shell.UpdateFile("file1", "first line\nsecond line\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1")).
			// Enlarge the layout while the files panel is focused. The main view
			// is only partially shown, so the renderer renders narrow columns.
			Press(keys.Universal.NextScreenMode)

		t.Views().Main().Content(Contains("COLUMNWIDTH=30"))

		// Focus the main view; it now takes up the full width. The renderer must
		// be re-run so that its columns use the full width.
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			/* EXPECTED:
			Content(Contains("COLUMNWIDTH=68"))
			ACTUAL: */
			Content(Contains("COLUMNWIDTH=30"))
	},
})
