package fingerprint

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

// CalcFingerprint returns the audio Fingerprint of the file at fPath.
// fPath can be a path to a directory or to a single file
func (c ChromaIO) CalcFingerprint(fPath string) ([]*Fingerprint, error) {
	fInfo, err := os.Stat(fPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrFileNotExists
		}

		return nil, err
	}

	if fInfo.IsDir() {
		return scanAudioDir(fPath)
	}

	if !isValidExtension(filepath.Ext(fInfo.Name())) {
		return nil, ErrInvalidFormat
	}

	fing, err := execFPcalc(fPath)
	if err != nil {
		return nil, err
	}

	return []*Fingerprint{fing}, nil
}

func scanAudioDir(dirPath string) ([]*Fingerprint, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var fings []*Fingerprint
	for _, f := range files {
		if f.IsDir() || !isValidExtension(filepath.Ext(f.Name())) {
			continue
		}

		fing, err := execFPcalc(path.Join(dirPath, f.Name()))
		if err != nil {
			return nil, err
		}

		fings = append(fings, fing)
	}

	return fings, nil
}

func execFPcalc(filePath string) (*Fingerprint, error) {
	cmd := exec.Command("fpcalc", "-json", filePath)
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
