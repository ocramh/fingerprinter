package types

import (
	"time"
)

// Release represents a real-world release objects (like a physical album) that can
// be bought in a music store.
// It is defined  according to https://musicbrainz.org/doc/Release
type Release struct {
	ID           string
	Title        string
	ReleasedAt   time.Time
	Format       string
	Tracks       []Track
	Author       []Author
	RecordLabels []Label
}

// Track is a release recording or unique audio data. It is loosely modelled
// according to https://musicbrainz.org/doc/Recording
type Track struct {
	DurationMillis int
	Title          string
	Position       int
	AudioURL       string
	ISCRCode       string
	Available      bool
}

// Author are the artists that the release is primarily credited to
type Author struct {
	ID   string
	Name string
}

// Label is the entity which issued the release
type Label struct {
	Name        string
	ID          string
	Description string
}
