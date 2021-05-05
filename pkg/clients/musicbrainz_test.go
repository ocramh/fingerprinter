package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var (
	testAppName    = "my-app"
	testAppVersion = "0.0.1"
	testEmail      = "foo@bar.com"
)

func TestNewMusicBrainzClient(t *testing.T) {
	got := NewMusicBrainz(testAppName, testAppVersion, testEmail)

	assert.Equal(t, &MusicBrainz{
		appName:      testAppName,
		appSemVer:    testAppVersion,
		contactEmail: testEmail,
	}, got)
}

func TestGetRecordingInfoOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var recordingID = "d4d24fa2-22f5-4b02-8751-8c0cf9cd02b2"

	params := url.Values{}
	params.Add("fmt", "json")
	params.Add("inc", strings.Join(RecordingInfoQueryVals, "+"))
	reqURL := fmt.Sprintf("https://musicbrainz.org/ws/2/recording/%s?%s", recordingID, params.Encode())

	testDataFilepath := "../../test/data/musicbrainz_recording.json"
	data, err := ioutil.ReadFile(testDataFilepath)
	assert.NoError(t, err)

	httpmock.RegisterResponder("GET", reqURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusOK, data)
			return resp, nil
		},
	)

	client := NewMusicBrainz(testAppName, testAppVersion, testEmail)

	got, err := client.GetRecordingInfo(recordingID)
	assert.NoError(t, err)
	assert.Equal(t, got.Title, "Dot Net")
}
