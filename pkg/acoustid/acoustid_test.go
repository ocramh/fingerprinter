package acoustid

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	hc "github.com/ocramh/fingerprinter/internal/httpclient"
	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
)

func TestLookupFingerprintOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testDataFilepath := "../../test/data/acoustid_response.json"
	data, err := ioutil.ReadFile(testDataFilepath)
	assert.NoError(t, err)

	httpmock.RegisterResponder("POST", AcoustIDBaseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusOK, data)
			return resp, nil
		},
	)

	acClient := NewAcoustID("secret-key")
	fingerprint := fp.Fingerprint{
		Duration: 100,
		Value:    "the-extracted-fingerprint",
	}
	got, err := acClient.LookupFingerprint(&fingerprint, false)
	assert.NoError(t, err)
	assert.Equal(t, &AcoustIDLookupResp{
		Status: "ok",
		Results: []ACLookupResult{
			{
				ID:    "033908fc-19da-4afa-a8a8-f8e1b87ada75",
				Score: 0.995636,
				Recordings: []Recording{
					{
						MBRecordingID: "d4d24fa2-22f5-4b02-8751-8c0cf9cd02b2",
						MBReleaseGroups: []ReleaseGroup{
							{
								ID:    "baca2dcc-b3e7-4e5f-9560-68513356125d",
								Title: "La Di Da Di",
								Type:  "Album",
								Releases: []Release{
									{ID: "c100950b-5000-402c-a0dc-eb334840d134"},
									{ID: "c0925a40-863c-4df7-bb22-8a2f74c124c2"},
									{ID: "6e1d42d8-0cd5-4774-8606-ce33687893bc"},
									{ID: "05ea68c9-0f99-4b18-bddc-3f3f584b6143"},
								},
							},
						},
					},
				},
			},
		},
	}, got)
}

func TestLookupFingerprintStatusNotOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testDataFilepath := "../../test/data/acoustid_err_response.json"
	data, err := ioutil.ReadFile(testDataFilepath)
	assert.NoError(t, err)

	httpmock.RegisterResponder("POST", AcoustIDBaseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusBadRequest, data)
			return resp, nil
		},
	)

	acClient := NewAcoustID("secret-key")
	fingerprint := fp.Fingerprint{
		Duration: 100,
		Value:    "the-extracted-fingerprint",
	}
	_, err = acClient.LookupFingerprint(&fingerprint, false)
	assert.Equal(t, hc.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "invalid fingerprint",
	}, err)
}

func TestLookupFingerprintStatusServiceUnavailable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", AcoustIDBaseURL,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(http.StatusServiceUnavailable, []byte{})
			return resp, nil
		},
	)

	acClient := NewAcoustID("secret-key")
	fingerprint := fp.Fingerprint{
		Duration: 100,
		Value:    "the-extracted-fingerprint",
	}
	_, err := acClient.LookupFingerprint(&fingerprint, true)
	assert.Equal(t, hc.HTTPError{
		Code:    http.StatusServiceUnavailable,
		Message: "upstream service not available",
	}, err)
	assert.Equal(t, 2, httpmock.GetTotalCallCount())
}
