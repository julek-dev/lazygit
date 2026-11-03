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

// GetFilesInDiff get the specified commit files
func (self *CommitFileLoader) GetFilesInDiff(from string, to string, reverse bool) ([]*models.CommitFile, error) {
	output, err := self.diff(from, to, reverse)
	if err != nil {
		return nil, err
	}

	return getCommitFilesFromRawDiff(output), nil
}

// diff lists the files of the diff between from and to.
func (self *CommitFileLoader) diff(from string, to string, reverse bool) (string, error) {
	cmdArgs := NewGitCmd("diff").
		ParseableDiffPrefixes().
		Arg("--submodule").
		Arg("--no-ext-diff").
		Arg("--raw").
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
// ":100644 100644 1234567 89abcde R100\x00old\x00new\x00".
func getCommitFilesFromRawDiff(output string) []*models.CommitFile {
	fields := strings.Split(output, "\x00")
	commitFiles := []*models.CommitFile{}
	for i := 0; i < len(fields) && strings.HasPrefix(fields[i], ":"); {
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
	}

	return commitFiles
}
