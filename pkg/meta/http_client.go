package meta

import (
	"net/http"
	"time"
)

const (
	// ReqTimeout is the POST resquets timeout in seconds
	ReqTimeout = 10 * time.Second
)

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: ReqTimeout,
	}
}
