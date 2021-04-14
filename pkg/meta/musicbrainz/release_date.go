package musicbrainz

import (
	"encoding/json"
	"time"
)

type ReleaseDate time.Time

func (m *ReleaseDate) UnmarshalJSON(b []byte) error {
	var releaseDateStr string
	if err := json.Unmarshal(b, &releaseDateStr); err != nil {
		return err
	}

	const shortForm = "2006-01-02"
	const xshortForm = "2006-01"
	parsed, err := time.Parse(shortForm, releaseDateStr)
	if err != nil {
		parsed, err = time.Parse(xshortForm, releaseDateStr)
		if err != nil {
			return err
		}
	}

	*m = ReleaseDate(parsed)
	return nil
}
