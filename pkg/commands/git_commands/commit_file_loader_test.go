package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestGetCommitFilesFromRawDiff(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		output   []*models.CommitFile
	}{
		{
			testName: "no files",
			input:    "",
			output:   []*models.CommitFile{},
		},
		{
			testName: "one file",
			input:    ":100644 100644 1234567 89abcde M\x00Myfile\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
			},
		},
		{
			testName: "several files",
			input: ":100644 100644 1234567 89abcde M\x00Myfile\x00" +
				":000000 100644 0000000 89abcde A\x00MyOtherFile\x00" +
				":100644 000000 1234567 0000000 D\x00YetAnother\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
				{
					Path:         "MyOtherFile",
					ChangeStatus: "A",
				},
				{
					Path:         "YetAnother",
					ChangeStatus: "D",
				},
			},
		},
		{
			testName: "a rename and a copy among regular files",
			input: ":100644 100644 1234567 89abcde M\x00Myfile\x00" +
				":100644 100644 1234567 1234567 R100\x00before\x00after\x00" +
				":100644 100644 1234567 89abcde C075\x00original\x00copy\x00" +
				":000000 100644 0000000 89abcde A\x00Added\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
				{
					Path:         "after",
					PreviousPath: "before",
					ChangeStatus: "R",
				},
				{
					Path:         "copy",
					PreviousPath: "original",
					ChangeStatus: "C",
				},
				{
					Path:         "Added",
					ChangeStatus: "A",
				},
			},
		},
		{
			testName: "a path that starts with a colon",
			input:    ":100644 100644 1234567 89abcde M\x00:file\x00",
			output: []*models.CommitFile{
				{
					Path:         ":file",
					ChangeStatus: "M",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result := getCommitFilesFromRawDiff(test.input)
			assert.Equal(t, test.output, result)
		})
	}
}
