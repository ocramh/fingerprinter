package verifier

import (
	"errors"
	"log"
	"sort"
	"time"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	meta "github.com/ocramh/fingerprinter/pkg/meta"
	mb_types "github.com/ocramh/fingerprinter/pkg/meta/musicbrainz"
)

type ReleaseGroupID string

// AudioVerifier is responsible for verifying the metadata integrity of individal
// audio files or folders
type AudioVerifier struct {
	chromaMngr     *fp.ChromaIO
	acClient       *meta.AcoustIDClient
	mbClient       *meta.MBClient
	acoustReleases map[ReleaseGroupID]meta.ReleaseGroup
}

func NewAudioVerifier(ch *fp.ChromaIO, ac *meta.AcoustIDClient, mb *meta.MBClient) *AudioVerifier {
	return &AudioVerifier{
		chromaMngr:     ch,
		acClient:       ac,
		mbClient:       mb,
		acoustReleases: make(map[ReleaseGroupID]meta.ReleaseGroup),
	}
}

func (a AudioVerifier) Analyze(inputPath string) (ra *RecAnalysis, err error) {
	fingerps, err := a.chromaMngr.CalcFingerprint(inputPath)
	if err != nil {
		return nil, err
	}

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
			log.Printf("[mb release-group ID] %s \n", recording.MBReleaseGroupsID)
			log.Printf("[mb recording ID] %s \n", recording.MBRecordingID)

			availableRecordingIDS = append(availableRecordingIDS, recording.MBRecordingID)

			for _, releaseGroup := range recording.MBReleaseGroupsID {
				_, ok := a.acoustReleases[ReleaseGroupID(releaseGroup.ID)]
				if !ok {
					a.acoustReleases[ReleaseGroupID(releaseGroup.ID)] = releaseGroup
				}
			}
		}
	}

	// remove duplicated recordings
	var analysis RecAnalysis
	for _, releaseGroupInfo := range a.acoustReleases {
		releaseData := ReleaseMeta{
			ID:         releaseGroupInfo.ID,
			Title:      releaseGroupInfo.Title,
			ReleasedAt: time.Now(),
			Authors:    []Author{},
			LabelInfo:  []Label{},
			Tracks:     []mb_types.Track{},
		}

		for _, release := range releaseGroupInfo.Releases {
			log.Printf("mb lookup release: %s \n", release.ID)
			releaseInfo, err := a.mbClient.GetReleaseInfo(string(release.ID))
			if err != nil {
				return nil, err
			}

			// set 1st release date
			if time.Time(releaseInfo.ReleasedAt).Before(releaseData.ReleasedAt) {
				releaseData.ReleasedAt = time.Time(releaseInfo.ReleasedAt)
			}

			for _, aut := range releaseInfo.Authors {
				if !releaseData.hasAuthor(aut) {
					releaseData.Authors = append(releaseData.Authors, Author{ID: aut.ArtistMeta.ID, Name: aut.Name, Description: aut.Description})
				}
			}

			for _, lab := range releaseInfo.LabelInfo {
				if !releaseData.hasLabel(lab.Label) {
					releaseData.LabelInfo = append(releaseData.LabelInfo, Label{ID: lab.Label.ID, Name: lab.Label.Name, Description: lab.Label.Description})
				}
			}

			for _, media := range releaseInfo.Media {
				for _, trk := range media.Tracks {
					if !releaseData.hasTrack(trk) {
						releaseData.Tracks = append(releaseData.Tracks, trk)
					}

					for _, availableRecID := range availableRecordingIDS {
						if trk.Recording.ID == availableRecID && !releaseData.hasAvailableTrack(availableRecID, trk.Recording.ISRCs) {
							releaseData.AvailableTracks = append(releaseData.AvailableTracks, trk)
						}
					}
				}
			}

			time.Sleep(meta.MusicBrainzReqDelay)
		}

		analysis.MatchedReleases = append(analysis.MatchedReleases, releaseData)
	}

	return &analysis, nil
}

type RecAnalysis struct {
	MatchedReleases []ReleaseMeta
}

type ReleaseMeta struct {
	ID              string
	Title           string
	ReleasedAt      time.Time
	Format          string
	Authors         []Author
	LabelInfo       []Label
	Tracks          []mb_types.Track
	AvailableTracks []mb_types.Track
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

func (r ReleaseMeta) hasAuthor(a mb_types.Author) bool {
	for _, aut := range r.Authors {
		if a.ArtistMeta.ID == aut.ID {
			return true
		}
	}
	return false
}

func (r ReleaseMeta) hasLabel(l mb_types.Label) bool {
	for _, lab := range r.LabelInfo {
		if l.ID == lab.ID {
			return true
		}
	}
	return false
}

func (r ReleaseMeta) hasTrack(rec mb_types.Track) bool {
	for _, trk := range r.Tracks {
		if trk.ID == rec.ID {
			return true
		}

		for _, knownIsrc := range trk.Recording.ISRCs {
			for _, newIsrc := range rec.Recording.ISRCs {
				if knownIsrc == newIsrc {
					return true
				}
			}
		}
	}

	return false
}

func (r ReleaseMeta) hasAvailableTrack(recID string, isrcs []string) bool {
	for _, rec := range r.AvailableTracks {
		if rec.ID == recID {
			return true
		}

		for _, newIsrc := range isrcs {
			for _, knownIsrc := range rec.Recording.ISRCs {
				if knownIsrc == newIsrc {
					return true
				}
			}
		}
	}

	return false
}
