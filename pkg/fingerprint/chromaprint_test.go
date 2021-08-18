package fingerprint

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDataDir = "../../test/data/audio/"
)

func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	// Print out the test value to stdout
	fmt.Fprintf(os.Stdout, `{"duration": 10.5, "fingerprint": "the-fingerprint"}`)
	os.Exit(0)
}

func TestFingerprintFromFile(t *testing.T) {
	mockExec := func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestShellProcessSuccess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_PROCESS=1"}
		return cmd
	}
	inputFilePath := path.Join(testDataDir, "sample1.mp3")
	chromap := NewChromaPrint(mockExec)
	fInfo, err := os.Stat(inputFilePath)
	assert.NoError(t, err)

	got, err := chromap.CalcFingerprint(inputFilePath)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, got[0], &Fingerprint{
		Duration:  10.5,
		Value:     "the-fingerprint",
		InputFile: fInfo,
	})
}

func TestFingerprintFromDir(t *testing.T) {
	assert.Equal(t, 1, 1)
}
