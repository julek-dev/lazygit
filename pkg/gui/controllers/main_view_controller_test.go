package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdjacentHunkHeaderIdx(t *testing.T) {
	lines := []string{
		"diff --git a/file1 b/file1",
		"index c4352f8..cc3b380 100644",
		"--- a/file1",
		"+++ b/file1",
		"@@ -1,4 +1,4 @@",
		"-line 1",
		"+CHANGED FIRST",
		" line 2",
		"@@ -17,4 +17,4 @@ line 16",
		" line 17",
		"-line 20",
		"+CHANGED LAST",
	}

	scenarios := []struct {
		name        string
		lines       []string
		fromIdx     int
		forward     bool
		expectedIdx int
		expectFound bool
	}{
		{"forward from top lands on first hunk", lines, 0, true, 4, true},
		{"forward from first hunk lands on second hunk", lines, 4, true, 8, true},
		{"forward past the last hunk finds nothing", lines, 8, true, 0, false},
		{"forward from within a hunk lands on the next one", lines, 6, true, 8, true},
		{"backward from the last hunk lands on the first", lines, 8, false, 4, true},
		{"backward from within a hunk lands on its header", lines, 6, false, 4, true},
		{"backward before the first hunk finds nothing", lines, 4, false, 0, false},
		{"no hunks at all", []string{" context", " context"}, 0, true, 0, false},
		{"empty input", []string{}, 0, true, 0, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			idx, found := adjacentHunkHeaderIdx(s.lines, s.fromIdx, s.forward)
			assert.Equal(t, s.expectFound, found)
			if s.expectFound {
				assert.Equal(t, s.expectedIdx, idx)
			}
		})
	}
}
