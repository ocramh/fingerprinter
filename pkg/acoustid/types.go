package acoustid

// AcoustIDLookupResp is the type used to parse a successfull AcoustID JSON response
type AcoustIDLookupResp struct {
	Status  string           `json:"status"`
	Results []ACLookupResult `json:"results"`
}

// ACLookupResult is a fingerprint match. It contaons one or more recordings that
// include the audio fingerprint analized and the accuracy score
type ACLookupResult struct {
	ID         string      `json:"id"`
	Score      float32     `json:"score"`
	Recordings []Recording `json:"recordings"`
}

// Recording is a single recording as defined by the MusicBrainz catalogue
type Recording struct {
	MBRecordingID   string         `json:"id"`
	MBReleaseGroups []ReleaseGroup `json:"releasegroups"`
}

// ReleaseGroup is a logical group of releases
type ReleaseGroup struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Type     string    `json:"type"`
	Releases []Release `json:"releases"`
}

// Release identifies a unique release on the MusicBrainz catalogue
type Release struct {
	ID string `json:"id"`
}

// AcoustErrResp is the type used to parse an AcoustID error JSON response
type AcoustErrResp struct {
	Error acoustIDErr `json:"error"`
}

type acoustIDErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
