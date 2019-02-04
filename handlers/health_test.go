package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emojify-app/api/emojify"
	"github.com/emojify-app/api/logging"
	"github.com/machinebox/sdk-go/boxutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupHealthTests() (*Health, *httptest.ResponseRecorder, *http.Request) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	l, _ := logging.New("test", "test", "localhost:8125", "DEBUG", "text")

	em := &emojify.MockEmojify{}
	em.On("Health", mock.Anything).Return(&boxutil.Info{}, nil)

	return &Health{l, em}, rw, r
}

func TestHealthHandlerReturnsOK(t *testing.T) {
	h, rw, r := setupHealthTests()

	h.ServeHTTP(rw, r)

	assert.Contains(t, string(rw.Body.Bytes()), "OK")
}

func TestHealthHandlerReturns200(t *testing.T) {
	h, rw, r := setupHealthTests()

	h.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
}