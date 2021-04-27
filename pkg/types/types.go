package types

import (
	"time"
)

type Release struct {
	Title        string
	ReleasedAt   time.Time
	Format       string
	Tracks       []Track
	Author       []Author
	RecordLabels []Label
}

type Track struct {
	DurationMillis int
	Title          string
	Position       int
	AudioURL       string
	ISCRCode       string
	Available      bool
}

type Author struct {
	ID   string
	Name string
}

type Label struct {
	Name        string
	ID          string
	Description string
}
