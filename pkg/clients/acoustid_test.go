package clients

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	"github.com/stretchr/testify/assert"
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
	got, err := acClient.LookupFingerprint(&fingerprint)
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
	_, err = acClient.LookupFingerprint(&fingerprint)
	assert.Equal(t, HTTPError{
		code:    http.StatusBadRequest,
		message: "invalid fingerprint",
	}, err)
}
