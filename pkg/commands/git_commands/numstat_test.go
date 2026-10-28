package git_commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNumstat(t *testing.T) {
	scenarios := []struct {
		name     string
		output   string
		expected map[string]FileDiff
	}{
		{
			name:     "no files",
			output:   "",
			expected: map[string]FileDiff{},
		},
		{
			name:   "files",
			output: "4\t1\tfile1\x000\t2\tdir/file2\x00",
			expected: map[string]FileDiff{
				"file1":     {LinesAdded: 4, LinesDeleted: 1},
				"dir/file2": {LinesDeleted: 2},
			},
		},
		{
			name:   "rename",
			output: "3\t1\t\x00before\x00after\x002\t0\tfile\x00",
			expected: map[string]FileDiff{
				"after": {LinesAdded: 3, LinesDeleted: 1, PreviousPath: "before"},
				"file":  {LinesAdded: 2},
			},
		},
		{
			name:   "binary files",
			output: "-\t-\timage.png\x00-\t-\t\x00old.png\x00new.png\x001\t1\tfile\x00",
			expected: map[string]FileDiff{
				"file": {LinesAdded: 1, LinesDeleted: 1},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, parseNumstat(strings.Split(s.output, "\x00")))
		})
	}
}
