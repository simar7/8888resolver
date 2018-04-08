package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dns "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetDNS_Type_A(t *testing.T) {
	defer gock.Off()

	testDNSRecord := DNSResponse{
		Answers: []DNSAnswer{
			{
				TTL:  42,
				Data: "2607:dead:face:cafe8::200e",
				Name: "example.com.",
				Type: dns.TypeAAAA,
			},
		},
		Questions: []DNSQuestion{
			{
				Name: "example.com.",
				Type: dns.TypeAAAA,
			},
		},
	}
	gockURL := fmt.Sprintf("%s?name=%s&type=%s", GoogleDNS, "example.com", "AAAA")
	gock.New(gockURL).Reply(http.StatusOK).JSON(testDNSRecord)

	ro := Router{
		HTTPClient: newHTTPClient(time.Second * 10),
	}
	router := ro.setupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dns/example.com/AAAA", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	expected := testDNSRecord
	actual := DNSResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected, actual)

	// Verify that we don't have pending mocks
	assert.True(t, gock.IsDone())
}
