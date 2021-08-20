package httpclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPErrorMessage(t *testing.T) {
	statusCode := http.StatusNotFound
	msg := "item not found"
	httperr := NewHTTPError(statusCode, msg)

	expected := fmt.Sprintf("http status: %d. error message: %s", statusCode, msg)
	assert.EqualError(t, httperr, expected)
}
