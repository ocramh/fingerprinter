package verifier

import (
	"errors"
	"log"
	"sort"
	"time"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	meta "github.com/ocramh/fingerprinter/pkg/meta"
)

type AcoustReleaseID string
type AcoustRecordingID string

// AudioVerifier is responsible for verifying the metadata integrity of individal
// audio files or folders
type AudioVerifier struct {
	chromaMngr     *fp.ChromaIO
	acClient       *meta.AcoustIDClient
	mbClient       *meta.MBClient
	acoustReleases map[AcoustReleaseID][]AcoustRecordingID
}

func NewAudioVerifier(ch *fp.ChromaIO, ac *meta.AcoustIDClient, mb *meta.MBClient) *AudioVerifier {
	return &AudioVerifier{
		chromaMngr:     ch,
		acClient:       ac,
		mbClient:       mb,
		acoustReleases: make(map[AcoustReleaseID][]AcoustRecordingID),
	}
}

func (a AudioVerifier) Analyze(inputPath string) (ra *RecAnalysis, err error) {
	fingerps, err := a.chromaMngr.CalcFingerprint(inputPath)
	if err != nil {
		return nil, err
	}

	log.Printf("generated %d fingerprints", len(fingerps))

	// query acoustid to match fingerprints with recordings (aka tracks) and get
	// associated releases (aka albums)
	var availableRecordingIDS []string
	for _, fingerp := range fingerps {
		acLookup, err := a.acClient.LookupFingerprint(fingerp)
		if err != nil {
			return nil, err
		}

		// order by score and get first one
		if len(acLookup.Results) == 0 {
			return nil, errors.New("no results found")
		}

		sort.Sort(meta.ACResultsByScore(acLookup.Results))
		topAcMatch := acLookup.Results[0]

		if len(topAcMatch.Recordings) == 0 {
			return nil, errors.New("no matches found on musicbrainz")
		}

		for _, recording := range topAcMatch.Recordings {
			log.Printf("[mb recording ID] %s \n", recording.MBRecordingID)
			recordingID := AcoustRecordingID(recording.MBRecordingID)
			availableRecordingIDS = append(availableRecordingIDS, recording.MBRecordingID)

			for _, releaseGroup := range recording.MBReleaseGroupsID {
				for _, release := range releaseGroup.Releases {
					releaseID := AcoustReleaseID(release.ID)
					_, ok := a.acoustReleases[releaseID]
					if ok {
						a.acoustReleases[releaseID] = append(a.acoustReleases[releaseID], recordingID)
					} else {
						a.acoustReleases[releaseID] = []AcoustRecordingID{recordingID}
					}
				}
			}
		}
	}

	// remove duplicated recordings
	var analysis RecAnalysis
	for releaseID := range a.acoustReleases {
		log.Printf("mb lookup release: %s \n", releaseID)
		releaseInfo, err := a.mbClient.GetReleaseInfo(string(releaseID))
		if err != nil {
			return nil, err
		}

		var authors []Author
		for _, aut := range releaseInfo.Authors {
			authors = append(authors, Author{ID: aut.ArtistMeta.ID, Name: aut.Name, Description: aut.Description})
		}

		var labels []Label
		for _, lab := range releaseInfo.LabelInfo {
			labels = append(labels, Label{ID: lab.Label.ID, Name: lab.Label.Name, Description: lab.Label.Description})
		}

		analysis.MatchedReleases = append(analysis.MatchedReleases, ReleaseMeta{
			ID:         releaseInfo.ID,
			Title:      releaseInfo.Title,
			ReleasedAt: time.Time(releaseInfo.ReleasedAt),
			Authors:    authors,
			LabelInfo:  labels,
		})

		for _, media := range releaseInfo.Media {
			for _, track := range media.Tracks {
				for _, availableRecID := range availableRecordingIDS {

					if len(track.Recording.ISRCs) > 0 {
						if track.Recording.ID == availableRecID && !analysis.hasAvailableRecording(availableRecID, track.Recording.ISRCs[0]) {
							analysis.AvailableRecordings = append(analysis.AvailableRecordings, Recording{
								ID:             track.ID,
								ISRC:           track.Recording.ISRCs[0],
								Title:          track.Title,
								DurationMillis: track.DurationMillis,
								Position:       track.Position,
							})
						}
					}
				}
			}
		}

		time.Sleep(meta.MusicBrainzReqDelay)
	}

	return &analysis, nil
}

type RecAnalysis struct {
	Valid               bool
	Complete            bool
	MatchedReleases     []ReleaseMeta
	AvailableRecordings []Recording
}

func (r RecAnalysis) hasAvailableRecording(recID, isrc string) bool {
	for _, rec := range r.AvailableRecordings {
		if rec.ID == recID || rec.ISRC == isrc {
			return true
		}
	}

	return false
}

type ReleaseMeta struct {
	ID         string
	Title      string
	ReleasedAt time.Time
	Format     string
	Authors    []Author
	LabelInfo  []Label
	Tracks     []Track
}

type Author struct {
	ID          string
	Name        string
	Description string
}

type Label struct {
	Name        string
	ID          string
	Description string
}

type Track struct{}

type Recording struct {
	ID             string
	ISRC           string
	Title          string
	DurationMillis int
	Position       int
}
