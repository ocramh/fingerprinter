package meta

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
)

const (
	// AcoustIDBaseURL is the base URL used for POSTing queries
	AcoustIDBaseURL = "https://api.acoustid.org/v2/lookup"

	// ReqTimeout is the POST resquets timeout in seconds
	ReqTimeout = 5 * time.Second
)

var (
	// lookupMeta is the metadata that will be added to a lookup response.
	// Recordings and releasegroups ids values can be used to query the MusicBrainz API
	lookupMeta = []string{"recordings", "recordingids", "releasegroups", "releasegroupids"}
)

// AcoustIDClient is the type responsible for interacting with the AcoustID API.
// It requires an API key that can be generated by registering an application at
// https://acoustid.org/login?return_url=https%3A%2F%2Facoustid.org%2Fnew-application
type AcoustIDClient struct {
	apiKey string
}

// NewAcoustIDClient is the AcoustIDClient cnstructor
func NewAcoustIDClient(k string) *AcoustIDClient {
	return &AcoustIDClient{k}
}

// LookupFingerprint uses audio fingerprints and duration values to search the AcoustID
// fingerprint database and return the corresponding track ID and MusicBrainz
// recording ID if a match was found
func (a *AcoustIDClient) LookupFingerprint(f *fp.Fingerprint) (*AcoustIDLookupResp, error) {
	encodedPayload := a.buildLookupQueryVals(f).Encode()
	req, err := http.NewRequest("POST", AcoustIDBaseURL, strings.NewReader(encodedPayload))
	if err != nil {
		return nil, err
	}

	// req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{
		Timeout: ReqTimeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrResp(resp)
	}

	var lookupResp AcoustIDLookupResp
	err = json.NewDecoder(resp.Body).Decode(&lookupResp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &lookupResp, nil
}

func (a *AcoustIDClient) buildLookupQueryVals(f *fp.Fingerprint) url.Values {
	values := url.Values{}
	values.Set("client", a.apiKey)
	values.Add("meta", strings.Join(lookupMeta, " "))
	values.Add("duration", strconv.Itoa(int(f.Duration)))
	values.Add("fingerprint", f.Value)

	return values
}

func handleErrResp(resp *http.Response) error {
	var errResp AcoustErrResp
	err := json.NewDecoder(resp.Body).Decode(&errResp)
	if err != nil {
		return err
	}

	return HTTPError{
		code:    errResp.Error.Code,
		message: errResp.Error.Message,
	}
}

// AcoustIDLookupResp is the type used to parse a successfull AcoustID JSON response
type AcoustIDLookupResp struct {
	Status  string         `json":"status"`
	Results []lookupResult `json:"results"`
}

type lookupResult struct {
	ID         string       `json:"id"`
	Recordings []recordings `json:"recordings"`
	Score      float32      `json:"score"`
}

type recordings struct {
	MBRecordingsID    string            `json:"id"`
	MBReleaseGroupsID []releaseGroupsID `json:"releasegroups"`
	Sources           int               `json:"sources"`
}

type releaseGroupsID struct {
	ID string `json:"id"`
}

// AcoustErrResp is the type used to parse an AcoustID error JSON response
type AcoustErrResp struct {
	Error acoustIDErr `json:"error"`
}

type acoustIDErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
