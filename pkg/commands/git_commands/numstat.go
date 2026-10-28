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
// their NULs, into the counts of each file by its path.
func parseNumstat(fields []string) map[string]FileDiff {
	fileDiffs := map[string]FileDiff{}
	for _, line := range fields {
		splitLine := strings.Split(line, "\t")
		if len(splitLine) != 3 {
			continue
		}

		linesAdded, err := strconv.Atoi(splitLine[0])
		if err != nil {
			continue
		}
		linesDeleted, err := strconv.Atoi(splitLine[1])
		if err != nil {
			continue
		}

		fileName := splitLine[2]
		fileDiffs[fileName] = FileDiff{
			LinesAdded:   linesAdded,
			LinesDeleted: linesDeleted,
		}
	}

	return fileDiffs
}
