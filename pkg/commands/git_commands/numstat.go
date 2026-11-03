package git_commands

import (
	"strconv"
	"strings"
)

type FileDiff struct {
	LinesAdded   int
	LinesDeleted int
}

// parseNumstat parses the records that `git diff --numstat -z` prints, split at
// their NULs, into the counts of each file by its path. A file that several
// records name gets the sum of their counts.
func parseNumstat(fields []string) map[string]FileDiff {
	fileDiffs := map[string]FileDiff{}
	for i := 0; i < len(fields); i++ {
		splitLine := strings.Split(fields[i], "\t")
		if len(splitLine) != 3 {
			continue
		}

		fileName := splitLine[2]
		// A rename has no path here; its old and new path follow as fields of
		// their own.
		if fileName == "" {
			if i+2 >= len(fields) {
				break
			}
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
			LinesAdded:   fileDiffs[fileName].LinesAdded + linesAdded,
			LinesDeleted: fileDiffs[fileName].LinesDeleted + linesDeleted,
		}
	}

	return fileDiffs
}
