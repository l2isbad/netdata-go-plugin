package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"net/http"
)

var (
	username    = "user"
	password    = "password"
	headerKey   = "X-Api-Key"
	headerValue = "secret"
)

func TestRawRequest_CreateRequest(t *testing.T) {
	rawRequest := RawRequest{
		Username: username,
		Password: password,
		Header: map[string]string{
			headerKey: headerValue,
		},
	}
	req, err := rawRequest.CreateHTTPRequest()
	assert.IsType(t, (*http.Request)(nil), req)

	user, pass, ok := req.BasicAuth()

	assert.Nil(t, err)
	assert.True(t, ok)
	assert.True(t, user == username && pass == password)
	assert.True(t, req.Header.Get(headerKey) == headerValue)
}
