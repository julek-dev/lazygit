package git_commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type CommitFileLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewCommitFileLoader(common *common.Common, cmd oscommands.ICmdObjBuilder) *CommitFileLoader {
	return &CommitFileLoader{
		Common: common,
		cmd:    cmd,
	}
}

// GetFilesInDiff get the specified commit files. withNumstat also counts the
// lines that the diff adds to and deletes from each of them.
func (self *CommitFileLoader) GetFilesInDiff(from string, to string, reverse bool, withNumstat bool) ([]*models.CommitFile, error) {
	output, err := self.diff(from, to, reverse, withNumstat)
	if err != nil && withNumstat {
		// Counting reads the content of every file, which can fail where listing
		// them doesn't, e.g. in a partial clone that is offline.
		self.Log.Error(err)
		output, err = self.diff(from, to, reverse, false)
	}
	if err != nil {
		return nil, err
	}

	return getCommitFilesFromRawDiff(output), nil
}

// diff lists the files of the diff between from and to, and with withNumstat,
// how many lines it changes in each of them.
func (self *CommitFileLoader) diff(from string, to string, reverse bool, withNumstat bool) (string, error) {
	cmdArgs := NewGitCmd("diff").
		ParseableDiffPrefixes().
		Arg("--submodule").
		Arg("--no-ext-diff").
		Arg("--raw").
		ArgIf(withNumstat, "--numstat").
		Arg("-z").
		Arg(fmt.Sprintf("--find-renames=%d%%", self.UserConfig().Git.RenameSimilarityThreshold)).
		ArgIf(reverse, "-R").
		Arg(from).
		Arg(to).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

// getCommitFilesFromRawDiff parses what `git diff --raw -z` prints: for each file, a
// record such as ":100644 100644 1234567 89abcde M\x00file\x00", whose last field
// is the file's status. A rename or a copy has a similarity score appended to its
// status, and two paths, the old one and the new one, e.g.
// ":100644 100644 1234567 89abcde R100\x00old\x00new\x00". The records of
// --numstat follow, if it was given.
func getCommitFilesFromRawDiff(output string) []*models.CommitFile {
	fields := strings.Split(output, "\x00")
	commitFiles := []*models.CommitFile{}
	submodulePaths := map[string]bool{}
	i := 0
	for i < len(fields) && strings.HasPrefix(fields[i], ":") {
		// The modes, the object names and the status
		meta := strings.Fields(fields[i][1:])
		if len(meta) != 5 {
			break
		}
		changeStatus := meta[4]

		if changeStatus[0] == 'R' || changeStatus[0] == 'C' {
			if i+2 >= len(fields) {
				break
			}
			// Drop the similarity score, so the rest of the code only has to deal
			// with a plain "R" or "C".
			commitFiles = append(commitFiles, &models.CommitFile{
				ChangeStatus: changeStatus[:1],
				PreviousPath: fields[i+1],
				Path:         fields[i+2],
			})
			i += 3
		} else {
			if i+1 >= len(fields) {
				break
			}
			commitFiles = append(commitFiles, &models.CommitFile{
				ChangeStatus: changeStatus,
				Path:         fields[i+1],
			})
			i += 2
		}

		// A submodule has mode 160000 on the side of the diff that it exists on.
		if meta[0] == "160000" || meta[1] == "160000" {
			submodulePaths[commitFiles[len(commitFiles)-1].Path] = true
		}
	}

	fileDiffs := parseNumstat(fields[i:])
	for _, file := range commitFiles {
		// For a submodule, git counts a line for the commit it was at and one
		// for the commit it is at, which says nothing about how much changed.
		if diff, ok := fileDiffs[file.Path]; ok && !submodulePaths[file.Path] {
			file.LinesAdded = diff.LinesAdded
			file.LinesDeleted = diff.LinesDeleted
		}
	}

	return commitFiles
}
