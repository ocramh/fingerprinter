package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	mb "github.com/ocramh/fingerprinter/pkg/meta/musicbrainz"
)

const (
	MusicBrainzRecordingURL = "https://musicbrainz.org/ws/2/recording"
	MusicBrainzReleaseURL   = "https://musicbrainz.org/ws/2/release"
	MusicBrainzReqDelay     = 1 * time.Second // MusicBrainz allows one request per second
)

// MusicBrainz is the type responsible for interacting with the MusicBrainz API.
// See https://musicbrainz.org/doc/MusicBrainz_AP for API docs
type MusicBrainz struct {
	appName      string
	appSemVer    string
	contactEmail string
}

// NewMusicBrainz is the MBHTTPClient constructor
func NewMusicBrainz(appName string, appSemVer string, email string) *MusicBrainz {
	return &MusicBrainz{
		appName:      appName,
		appSemVer:    appSemVer,
		contactEmail: email,
	}
}

// GetRecordingInfo returns a single recording (or track) metadata.
// Metadata includes ISRC codes, releases info, recording titie, duration,
// release date, artists etc
func (m *MusicBrainz) GetRecordingInfo(recordingID string) (*mb.RecordingInfo, error) {
	includeVals := []string{"artists", "isrcs", "releases"}
	req, err := m.newMBGETRequest(MusicBrainzRecordingURL, recordingID, includeVals)
	if err != nil {
		return nil, err
	}

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, m.handleMBErrResp(resp)
	}

	var recInfo mb.RecordingInfo
	err = json.NewDecoder(resp.Body).Decode(&recInfo)
	if err != nil {
		return nil, err
	}

	return &recInfo, nil
}

// GetReleaseInfo returns a release metadata. Releases a real-world release objects
// such as a physical album that contains one or more Recordings
func (m *MusicBrainz) GetReleaseInfo(releaseID string) (*mb.ReleaseInfo, error) {
	includeVals := []string{"artists", "labels", "isrcs", "recordings", "artist-credits"}
	req, err := m.newMBGETRequest(MusicBrainzReleaseURL, releaseID, includeVals)
	if err != nil {
		return nil, err
	}

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, m.handleMBErrResp(resp)
	}

	var relInfo mb.ReleaseInfo
	err = json.NewDecoder(resp.Body).Decode(&relInfo)
	if err != nil {
		return nil, err
	}

	return &relInfo, nil
}

func (m *MusicBrainz) handleMBErrResp(r *http.Response) error {
	var errResp mb.MBError
	err := json.NewDecoder(r.Body).Decode(&errResp)
	if err != nil {
		return err
	}

	return HTTPError{
		code:    r.StatusCode,
		message: errResp.Error,
	}
}

// newMBGETRequest builds a new MusicBrainz HTTP GET request.
// It takes care of setting the right headers and url formatting
func (m *MusicBrainz) newMBGETRequest(baseURL string, entityID string, inc []string) (*http.Request, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, entityID))
	if err != nil {
		return nil, err
	}

	reqParams := url.Values{}
	reqParams.Set("fmt", "json")
	reqParams.Add("inc", strings.Join(inc, "+"))

	reqURL.RawQuery = reqParams.Encode()

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// see https://musicbrainz.org/doc/MusicBrainz_API/Rate_Limiting#Provide_meaningful_User-Agent_strings
	userAgent := fmt.Sprintf("%s/%s ( %s )", m.appName, m.appSemVer, m.contactEmail)
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}
