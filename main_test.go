package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDNS_Type_A(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dns/example.com/A", nil)
	router.ServeHTTP(w, req)

	actual := DNSResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)

	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, actual)
}
