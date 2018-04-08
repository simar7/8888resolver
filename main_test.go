package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dns "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetDNS(t *testing.T) {
	defer gock.Off()

	testCases := []struct {
		Name             string
		TestDNSRecord    DNSResponse
		ExpectedError    error
		ExpectedResponse DNSResponse
		DNSURL           string
		RequestURL       string
		HTTPStatusCode   int
	}{
		{
			Name: "Get Type AAAA",
			TestDNSRecord: DNSResponse{
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
			},
			ExpectedError: nil,
			ExpectedResponse: DNSResponse{
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
			},
			DNSURL:         fmt.Sprintf("%s?name=%s&type=%s", GoogleDNS, "example.com", "AAAA"),
			RequestURL:     fmt.Sprintf("/dns/example.com/AAAA"),
			HTTPStatusCode: http.StatusOK,
		},
		{
			Name: "Get Type A",
			TestDNSRecord: DNSResponse{
				Answers: []DNSAnswer{
					{
						TTL:  42,
						Data: "8.8.8.8",
						Name: "example.com.",
						Type: dns.TypeA,
					},
				},
				Questions: []DNSQuestion{
					{
						Name: "example.com.",
						Type: dns.TypeA,
					},
				},
			},
			ExpectedError: nil,
			ExpectedResponse: DNSResponse{
				Answers: []DNSAnswer{
					{
						TTL:  42,
						Data: "8.8.8.8",
						Name: "example.com.",
						Type: dns.TypeA,
					},
				},
				Questions: []DNSQuestion{
					{
						Name: "example.com.",
						Type: dns.TypeA,
					},
				},
			},
			DNSURL:         fmt.Sprintf("%s?name=%s&type=%s", GoogleDNS, "example.com", "A"),
			RequestURL:     fmt.Sprintf("/dns/example.com/A"),
			HTTPStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		gock.New(tc.DNSURL).Reply(tc.HTTPStatusCode).JSON(tc.TestDNSRecord)

		ro := Router{
			HTTPClient: newHTTPClient(time.Second * 10),
		}
		router := ro.setupRouter()

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, tc.RequestURL, nil)
		assert.Equal(t, tc.ExpectedError, err, tc.Name)

		router.ServeHTTP(w, req)

		expected := tc.ExpectedResponse
		actual := DNSResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &actual)
		assert.NoError(t, err, tc.Name)

		assert.Equal(t, tc.HTTPStatusCode, w.Code, tc.Name)
		assert.Equal(t, expected, actual, tc.Name)
	}

	// Verify that we don't have pending mocks
	assert.True(t, gock.IsDone())
}

func TestGetDNS_ErrorReachingResolver(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("%s?name=%s&type=%s", GoogleDNS, "bad.url", "A")).ReplyError(errors.New("internal server error"))
	ro := Router{
		HTTPClient: newHTTPClient(time.Second * 10),
	}
	router := ro.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/dns/bad.url/A"), nil)

	router.ServeHTTP(w, req)

	expected := DNSError{Error: "internal server error"}
	actual := DNSError{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expected, actual)
}
