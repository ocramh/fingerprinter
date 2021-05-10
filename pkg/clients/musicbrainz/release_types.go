package musicbrainz

// ReleaseInfo is a release info response type returned by the MusicBrainz API
type ReleaseInfo struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	LabelInfo  []LabelInfo `json:"label-info"`
	Authors    []Author    `json:"artist-credit"`
	Media      []Media     `json:"media"`
	ReleasedAt ReleaseDate `json:"date"`
}

type Media struct {
	Format string  `json:"format"`
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Title          string    `json:"title"`
	DurationMillis int       `json:"length"`
	Position       int       `json:"position"`
	ID             string    `json:"id"`
	Recording      Recording `json:"recording"`
	Authors        []Author  `json:"artist-credit"`
}

type Recording struct {
	ISRCs []string `json:"isrcs"`
	ID    string   `json:"id"`
}

type Author struct {
	Name        string     `json:"name"`
	ArtistMeta  ArtistMeta `json:"artist"`
	Description string     `json:"disambiguation"`
}

type ArtistMeta struct {
	ID string `json:"id"`
}

type LabelInfo struct {
	Label Label `json:"label"`
}

type Label struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"disambiguation"`
}
