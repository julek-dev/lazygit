package git_commands

import (
	"strconv"
	"strings"
)

type FileDiff struct {
	LinesAdded   int
	LinesDeleted int

	// The path a rename comes from, or "" if the diff doesn't see a rename.
	PreviousPath string
}

// parseNumstat parses the records that `git diff --numstat -z` prints, split at
// their NULs, into the counts of each file by its path.
func parseNumstat(fields []string) map[string]FileDiff {
	fileDiffs := map[string]FileDiff{}
	for i := 0; i < len(fields); i++ {
		splitLine := strings.Split(fields[i], "\t")
		if len(splitLine) != 3 {
			continue
		}

		fileName := splitLine[2]
		previousPath := ""
		// A rename has no path here; its old and new path follow as fields of
		// their own.
		if fileName == "" {
			if i+2 >= len(fields) {
				break
			}
			previousPath = fields[i+1]
			fileName = fields[i+2]
			i += 2
		}

		linesAdded, err := strconv.Atoi(splitLine[0])
		if err != nil {
			continue
		}
		linesDeleted, err := strconv.Atoi(splitLine[1])
		if err != nil {
			continue
		}

		fileDiffs[fileName] = FileDiff{
			LinesAdded:   linesAdded,
			LinesDeleted: linesDeleted,
			PreviousPath: previousPath,
		}
	}

	return fileDiffs
}
