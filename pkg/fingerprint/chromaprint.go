package fingerprint

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
)

var (
	ValidAudioFormats = []string{".mp3"}
)

// ChromaPrint is a concrete implementation of the Fingerprinter interface.
// It manages IO operations using the chromaprint library, which must be installed
// as external dependency
type ChromaPrint struct{}

// CalcFingerprint returns the audio Fingerprint of the file at fPath.
// fPath can be a path to a directory or to a single file
func (c ChromaPrint) CalcFingerprint(fPath string) ([]*Fingerprint, error) {
	fInfo, err := os.Stat(fPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrInvalidPath
		}

		return nil, err
	}

	if fInfo.IsDir() {
		return scanAudioDir(fPath)
	}

	return fingerprintFromFile(fInfo, fPath)
}

// result is the product of reading a file and extracting its adio fingerprint
type result struct {
	path   string
	fprint *Fingerprint
	err    error
}

// scanAudioDir scans the directory at dirPath and concurrently extracts fingerprints
func scanAudioDir(dirPath string) ([]*Fingerprint, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	fChan := make(chan result)
	fings := []*Fingerprint{}

	for _, fInfo := range files {
		if fInfo.IsDir() || !isValidExtension(filepath.Ext(fInfo.Name())) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			fpath := path.Join(dirPath, fInfo.Name())
			fing, err := execFPcalc(fInfo, fpath)
			fChan <- result{
				path:   fpath,
				fprint: fing,
				err:    err,
			}
		}()
	}

	go func() {
		for result := range fChan {
			if result.err != nil {
				fings = append(fings, result.fprint)
			}
		}
	}()

	wg.Wait()
	close(fChan)

	return fings, nil
}

func fingerprintFromFile(fInfo os.FileInfo, fPath string) ([]*Fingerprint, error) {
	if !isValidExtension(filepath.Ext(fInfo.Name())) {
		return nil, ErrInvalidFormat
	}

	fing, err := execFPcalc(fInfo, fPath)
	if err != nil {
		return nil, err
	}

	return []*Fingerprint{fing}, nil
}

func execFPcalc(fInfo os.FileInfo, fPath string) (*Fingerprint, error) {
	fpcalcExecPath, err := exec.LookPath("fpcalc")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(fpcalcExecPath, "-json", fPath)
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

	fp.InputFile = fInfo

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
