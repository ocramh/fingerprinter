package types

// RecordingInfo is a recording info response returned by the MusciBrainz API
type RecordingInfo struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Isrcs            []string       `json:"isrcs"`
	DurationMillisec int            `json:"length"`
	Releases         []releasesInfo `json:"releases"`
	ArtistCredit     []artistInfo   `json:"artist-credit"`
	ReleasedAt       ReleaseDate    `json:"first-release-date"`
}

type releasesInfo struct {
	Title   string `json:"title"`
	Country string `json:"country"`
	ID      string `json:"id"`
}

type artistInfo struct {
	Name        string     `json:"name"`
	ArtistMeta  artistMeta `json:"artist"`
	Description string     `json:"disambiguation"`
}

type artistMeta struct {
	ID string `json:"id"`
}
