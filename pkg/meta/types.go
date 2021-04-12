package meta

import (
	"time"

	mb "github.com/ocramh/fingerprinter/pkg/meta/musicbrainz"
)

type ReleaseData struct {
	Title        string
	ReleasedAt   time.Time
	Format       string
	Tracks       []TrackData
	Author       []Artist
	RecordLabels []mb.LabelInfo
}

type TrackData struct {
	DurationMillis int
	Title          string
	Position       int
	AudioURL       string
	ISCR           string
	Available      bool
}

type Artist struct {
	ID   string
	Name string
}
