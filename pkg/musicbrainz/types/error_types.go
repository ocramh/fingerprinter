package types

// MBError is a MusicBrainz API error response type
type MBError struct {
	Error string `json:"error"`
	Help  string `json:"help"`
}
