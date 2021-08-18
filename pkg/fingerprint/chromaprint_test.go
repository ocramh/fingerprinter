package fingerprint

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	testDataDir = "/test/audio/"
	testFile1   = "sample1.mp3"
	testFile2   = "sample2.mp3"
	testFile3   = "textfile.txt"
)

var (
	mockExec = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestShellProcessSuccess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_PROCESS=1"}
		return cmd
	}

	mockFailExec = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestShellProcessError", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_PROCESS=1"}
		return cmd
	}
)

func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	// Print out the test value to stdout
	fmt.Fprintf(os.Stdout, `{"duration": 10.5, "fingerprint": "the-fingerprint"}`)
	os.Exit(0)
}

func TestShellProcessError(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	os.Exit(2)
}

func mustSetupFS() afero.Fs {
	mockFS := afero.NewMemMapFs()
	err := mockFS.MkdirAll(testDataDir, 0755)
	if err != nil {
		panic(err)
	}

	inputFile1Path := path.Join(testDataDir, testFile1)
	err = afero.WriteFile(mockFS, inputFile1Path, []byte("file 1"), 0644)
	if err != nil {
		panic(err)
	}

	inputFile2Path := path.Join(testDataDir, testFile2)
	err = afero.WriteFile(mockFS, inputFile2Path, []byte("file 2"), 0644)
	if err != nil {
		panic(err)
	}

	inputFile3Path := path.Join(testDataDir, testFile3)
	err = afero.WriteFile(mockFS, inputFile3Path, []byte("text file"), 0644)
	if err != nil {
		panic(err)
	}

	return mockFS
}

func TestFingerprintFromFile(t *testing.T) {
	mockFS := mustSetupFS()

	chromap := NewChromaPrint(mockExec, mockFS)
	inputFile := path.Join(testDataDir, testFile1)
	fInfo, err := mockFS.Stat(inputFile)
	assert.NoError(t, err)

	got, err := chromap.CalcFingerprint(inputFile)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, got[0], &Fingerprint{
		Duration:  10.5,
		Value:     "the-fingerprint",
		InputFile: fInfo,
	})
}

func TestFingerprintFromDir(t *testing.T) {
	mockFS := mustSetupFS()

	chromap := NewChromaPrint(mockExec, mockFS)

	got, err := chromap.CalcFingerprint(testDataDir)
	assert.NoError(t, err)
	assert.Len(t, got, 2)

	fInfo1, err := mockFS.Stat(path.Join(testDataDir, testFile1))
	assert.NoError(t, err)

	fInfo2, err := mockFS.Stat(path.Join(testDataDir, testFile2))
	assert.NoError(t, err)

	assert.ElementsMatch(t, []*Fingerprint{
		{
			Duration:  10.5,
			Value:     "the-fingerprint",
			InputFile: fInfo1,
		},
		{
			Duration:  10.5,
			Value:     "the-fingerprint",
			InputFile: fInfo2,
		},
	}, got)
}

func TestInputErrors(t *testing.T) {
	mockFS := mustSetupFS()

	chromap := NewChromaPrint(mockExec, mockFS)

	testcases := []struct {
		name        string
		inputPath   string
		expectedErr error
	}{
		{
			name:        "invalid path",
			inputPath:   "some/other/dir",
			expectedErr: ErrInvalidPath,
		},
		{
			name:        "invalid file format",
			inputPath:   path.Join(testDataDir, testFile3),
			expectedErr: ErrInvalidFormat,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			_, err := chromap.CalcFingerprint(testcase.inputPath)
			assert.EqualError(t, err, testcase.expectedErr.Error())
		})
	}
}

func TestHandleExecCmdError(t *testing.T) {
	mockFS := mustSetupFS()

	chromap := NewChromaPrint(mockFailExec, mockFS)
	_, err := chromap.CalcFingerprint(testDataDir)
	assert.NotNil(t, err)
}
