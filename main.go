package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DNSQuestion is a DNS question struct with Name and Q-Type
type DNSQuestion struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
}

// DNSAnswer is a DNS answer struct with Name, Q-Type, TTL and Data
type DNSAnswer struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// DNSResponse amalgamates Questions and Answers slices
type DNSResponse struct {
	Questions []DNSQuestion `json:"Question"`
	Answers   []DNSAnswer   `json:"Answer"`
}

// Router is the HTTP router
type Router struct {
	HTTPClient *http.Client
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// GetDNS returns a DNSResponse of Questions and Answers
func (ro Router) GetDNS(c *gin.Context) {
	domain := c.Param("domain")
	qType := c.Param("qtype")

	dnsReq := fmt.Sprintf("%s?name=%s&type=%s", GoogleDNS, domain, qType)
	log.Println("Requesting: ", dnsReq)
	resp, err := ro.HTTPClient.Get(dnsReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}

	dnsResp := DNSResponse{}
	err = json.Unmarshal(body, &dnsResp)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}

	// return answer back
	c.JSON(http.StatusOK, dnsResp)
}

func (ro Router) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/dns/:domain/:qtype", ro.GetDNS)

	return r
}

func main() {
	ro := Router{
		HTTPClient: newHTTPClient(time.Second * 10),
	}
	r := ro.setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}
