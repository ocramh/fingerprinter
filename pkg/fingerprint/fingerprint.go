package fingerprint

import "os"

// Fingerprinter defines operations for calculating fingerprints from audio files
type Fingerprinter interface {

	// CalcFingerprint returns a list of fingerprints from an input path
	CalcFingerprint(fPath string) ([]*Fingerprint, error)
}

// Fingerprint is an audio file fingerprint. The JSON structure allows the struct to
// parse the chromaprint fpcalc command when executed with the -json flag
type Fingerprint struct {
	Duration  float32     `json:"duration"`
	Value     string      `json:"fingerprint"`
	InputFile os.FileInfo `json:"-"`
}
