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

const GoogleDNS string = "https://dns.google.com/resolve"

type DNSQuestion struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
}

type DNSAnswer struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

type DNSResponse struct {
	Questions []DNSQuestion `json:"Question"`
	Answers   []DNSAnswer   `json:"Answer"`
}

type Router struct {
	HTTPClient *http.Client
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func (ro Router) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/dns/:domain/:qtype", func(c *gin.Context) {
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
	})
	return r
}

func main() {
	ro := Router{
		HTTPClient: newHTTPClient(time.Second * 10),
	}
	r := ro.setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}
