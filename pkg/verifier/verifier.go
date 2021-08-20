package verifier

import (
	"log"
	"path"
	"sort"
	"time"

	ac "github.com/ocramh/fingerprinter/pkg/acoustid"
	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	mb "github.com/ocramh/fingerprinter/pkg/musicbrainz"
	mb_types "github.com/ocramh/fingerprinter/pkg/musicbrainz/types"
)

type ReleaseGroupID string

// AudioVerifier is responsible for verifying the metadata integrity of individal
// audio files or folders
type AudioVerifier struct {
	fprinter       fp.Fingerprinter
	acClient       *ac.AcoustID
	mbClient       *mb.MusicBrainz
	acoustReleases map[ReleaseGroupID]ac.ReleaseGroup
}

// AvailableRecording contains the uploaded file path and its associated musicbrainz
// recording ID
type AvailableRecording struct {
	ID       string
	FilePath string
}

func NewAudioVerifier(fp fp.Fingerprinter, acID *ac.AcoustID, mb *mb.MusicBrainz) *AudioVerifier {
	return &AudioVerifier{
		fprinter:       fp,
		acClient:       acID,
		mbClient:       mb,
		acoustReleases: make(map[ReleaseGroupID]ac.ReleaseGroup),
	}
}

func (a AudioVerifier) Analyze(inputPath string) (ra *RecAnalysis, err error) {
	fingerps, err := a.fprinter.CalcFingerprint(inputPath)
	if err != nil {
		return nil, err
	}

	// query acoustid to match fingerprints with recordings (aka tracks) and get
	// associated releases (aka albums)
	var availableRecordings []AvailableRecording
	var unmatchedAudioFiles []UnmatchedFile
	var retryOnFail = true
	for _, fingerp := range fingerps {
		acLookup, err := a.acClient.LookupFingerprint(fingerp, retryOnFail)
		if err != nil {
			return nil, err
		}

		// order by score and get first one
		if len(acLookup.Results) == 0 {
			log.Printf("no results found for %s", fingerp.InputFile.Name())
			unmatchedAudioFiles = append(unmatchedAudioFiles, UnmatchedFile{
				FileName: fingerp.InputFile.Name(),
				Reason:   "audio file fingerprint didn't match any known record",
			})
			continue
		}

		sort.Sort(ac.ACResultsByScore(acLookup.Results))
		topAcMatch := acLookup.Results[0]

		if len(topAcMatch.Recordings) == 0 {
			log.Printf("no recordings found for %s", fingerp.InputFile.Name())
			unmatchedAudioFiles = append(unmatchedAudioFiles, UnmatchedFile{
				FileName: fingerp.InputFile.Name(),
				Reason:   "audio file fingerprint didn't match any known release",
			})
		}

		for _, recording := range topAcMatch.Recordings {
			log.Printf("[mb recording ID] %s \n", recording.MBRecordingID)

			availableRecordings = append(availableRecordings, AvailableRecording{recording.MBRecordingID, path.Join(inputPath, fingerp.InputFile.Name())})

			for _, releaseGroup := range recording.MBReleaseGroups {
				releaseGroupInfo, ok := a.acoustReleases[ReleaseGroupID(releaseGroup.ID)]
				if !ok {
					a.acoustReleases[ReleaseGroupID(releaseGroup.ID)] = releaseGroup
				} else {
					a.acoustReleases[ReleaseGroupID(releaseGroup.ID)] = *addMissingReleasesIDToGroup(&releaseGroup, &releaseGroupInfo)
				}
			}
		}
	}

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

			releaseInfo, err := a.mbClient.GetReleaseInfo(release.ID)
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
					if len(trk.Recording.ISRCs) == 0 {
						continue
					}

					if !releaseData.hasTrack(trk) {
						releaseData.Tracks = append(releaseData.Tracks, trk)
					}

					for _, availableRec := range availableRecordings {
						if trk.Recording.ID == availableRec.ID && !releaseData.hasAvailableTrack(availableRec.ID, trk.Recording.ISRCs) {
							releaseData.AvailableTracks = append(releaseData.AvailableTracks, AvailableTrack{
								Track: trk,
								Path:  availableRec.FilePath,
							})
						}
					}
				}
			}

			time.Sleep(mb.MusicBrainzReqDelay)
		}

		analysis.MatchedReleases = append(analysis.MatchedReleases, releaseData)
	}
	analysis.UnmatchedFiles = unmatchedAudioFiles

	return &analysis, nil
}

// adds new releases IDs to the existing ReleaseGroup if they
// if they are not already included
func addMissingReleasesIDToGroup(new *ac.ReleaseGroup, existing *ac.ReleaseGroup) *ac.ReleaseGroup {
	for _, rg1 := range new.Releases {

		var releaseAlreadyExists bool
		for _, rg2 := range existing.Releases {
			if rg1.ID == rg2.ID {
				releaseAlreadyExists = true
			}
		}

		if !releaseAlreadyExists {
			existing.Releases = append(existing.Releases, rg1)
		}
	}

	return existing
}

// RecAnalysis is the result of performing audio analysis on a group of files
type RecAnalysis struct {
	MatchedReleases []ReleaseMeta
	UnmatchedFiles  []UnmatchedFile
}

// ReleaseMeta contains metadata that describes a single release
type ReleaseMeta struct {
	ID              string
	Title           string
	ReleasedAt      time.Time
	Format          string
	Authors         []Author
	LabelInfo       []Label
	Tracks          []mb_types.Track
	AvailableTracks []AvailableTrack
}

type AvailableTrack struct {
	Track mb_types.Track
	Path  string
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

// UnmatchedFile descirbes an audio file that doesn't match any entry on either
// acoustID or musicbrainz
type UnmatchedFile struct {
	FileName string
	Reason   string
}

func (r ReleaseMeta) hasAvailableTrack(recID string, isrcs []string) bool {
	for _, rec := range r.AvailableTracks {
		if rec.Track.ID == recID {
			return true
		}

		for _, newIsrc := range isrcs {
			for _, knownIsrc := range rec.Track.Recording.ISRCs {
				if knownIsrc == newIsrc {
					return true
				}
			}
		}
	}

	return false
}
