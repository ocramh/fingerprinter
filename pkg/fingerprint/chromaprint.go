package fingerprint

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
)

var (
	ValidAudioFormats = []string{".mp3"}
)

type ExecCmd = func(name string, arg ...string) *exec.Cmd

// ChromaPrint is a concrete implementation of the Fingerprinter interface.
// It manages IO operations using the chromaprint library, which must be installed
// as external dependency
type ChromaPrint struct {
	execCmd ExecCmd
	os      afero.Fs
}

func NewChromaPrint(exec ExecCmd, os afero.Fs) *ChromaPrint {
	return &ChromaPrint{
		execCmd: exec,
		os:      os,
	}
}

// CalcFingerprint returns the audio Fingerprint of the file at fPath.
// fPath can be a path to a directory or to a single file
func (c *ChromaPrint) CalcFingerprint(fPath string) ([]*Fingerprint, error) {
	fInfo, err := c.fileinfoFromPath(fPath)
	if err != nil {
		return nil, err
	}

	if fInfo.IsDir() {
		return c.scanAudioDir(fPath)
	}

	return c.fingerprintFromFile(fInfo, fPath)
}

// result is the product of reading a file and extracting its adio fingerprint
type result struct {
	path   string
	fprint *Fingerprint
	err    error
}

// scanAudioDir scans the directory at dirPath and concurrently extracts fingerprints.
// Subdirectories and files with an invalid extension will be ignored.
// This function only scans the top level directory
func (c *ChromaPrint) scanAudioDir(dirPath string) ([]*Fingerprint, error) {
	files, err := afero.ReadDir(c.os, dirPath)
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
		go func(fInfo os.FileInfo) {
			defer wg.Done()

			fpath := path.Join(dirPath, fInfo.Name())
			fing, err := c.execFPcalc(fInfo, fpath)

			fChan <- result{
				path:   fpath,
				fprint: fing,
				err:    err,
			}
		}(fInfo)
	}

	go func() {
		wg.Wait()
		close(fChan)
	}()

	// collect results
	for result := range fChan {
		if result.err != nil {
			return nil, result.err
		}

		fings = append(fings, result.fprint)
	}

	return fings, nil
}

func (c *ChromaPrint) fingerprintFromFile(fInfo os.FileInfo, fPath string) ([]*Fingerprint, error) {
	if !isValidExtension(filepath.Ext(fInfo.Name())) {
		return nil, ErrInvalidFormat
	}

	fing, err := c.execFPcalc(fInfo, fPath)
	if err != nil {
		return nil, err
	}

	return []*Fingerprint{fing}, nil
}

func (c *ChromaPrint) execFPcalc(fInfo os.FileInfo, fPath string) (*Fingerprint, error) {
	fpcalcExecPath, err := exec.LookPath("fpcalc")
	if err != nil {
		return nil, err
	}

	cmd := c.execCmd(fpcalcExecPath, "-json", fPath)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	var fp Fingerprint
	if err := json.NewDecoder(buf).Decode(&fp); err != nil {
		return nil, err
	}

	fp.InputFile = fInfo

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

func (c *ChromaPrint) fileinfoFromPath(p string) (os.FileInfo, error) {
	fInfo, err := c.os.Stat(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrInvalidPath
		}

		return nil, err
	}

	return fInfo, nil
}
