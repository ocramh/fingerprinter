package fingerprint

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ValidAudioFormats = []string{".mp3"}
)

// ChromaIO manages IO operations with the chromaprint library, which must be
// installed and availble as external dependency
type ChromaIO struct{}

// Fingerprint is the result of an audio file fingerprint. The JSON structure allows
// the struct to parse the chromaprint fpcalc command when executed with the -json
// flag
type Fingerprint struct {
	Duration float32 `json:"duration"`
	Value    string  `json:"fingerprint"`
}

// CalcFingerprint returns the audio Fingerprint of the file at fPath
func (c ChromaIO) CalcFingerprint(fPath string) (*Fingerprint, error) {
	fInfo, err := os.Stat(fPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrFileNotExists
		}

		return nil, err
	}

	if fInfo.IsDir() {
		return nil, ErrInvalidFileInput
	}

	if !isValidExtension(filepath.Ext(fInfo.Name())) {
		return nil, ErrInvalidFormat
	}

	cmd := exec.Command("fpcalc", "-json", fPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var fp Fingerprint
	if err := json.NewDecoder(stdout).Decode(&fp); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return &fp, nil
}

func isValidExtension(ext string) bool {
	for _, e := range ValidAudioFormats {
		if e == ext {
			return true
		}
	}

	return false
}
