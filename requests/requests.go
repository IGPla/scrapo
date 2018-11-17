package requests

import (
	"bytes"
	"github.com/IGPla/scrapo/config"
	"github.com/IGPla/scrapo/logger"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var requestsLogger *log.Logger

func init() {
	requestsLogger = logger.GetLogger("REQUESTS", os.Stdout)
}

/* Base function to get resource by url */
func GetResource(url string) (*bytes.Buffer, int, http.Header, error) {
	var result *bytes.Buffer = new(bytes.Buffer)

	var client *http.Client = buildClient()

	request, reqError := buildRequest(client, url)
	if reqError != nil {
		return nil, -1, nil, reqError
	}

	response, respError := getResponse(client, request)
	if respError != nil {
		return nil, -1, nil, respError
	}

	defer response.Body.Close()

	_, copyError := io.Copy(result, response.Body)
	if copyError != nil {
		return nil, -1, nil, copyError
	}

	return result, response.StatusCode, response.Header, nil
}

/* Build base client for http requests */
func buildClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

/* Build request */
func buildRequest(client *http.Client, url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		requestsLogger.Printf("Error arised creating request for %v (%v)",
			url,
			err.Error())
		return nil, err
	}
	request.Header.Set("User-Agent", config.MainConfig.UserAgent)
	return request, nil
}

/* Get response */
func getResponse(client *http.Client, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		requestsLogger.Printf("Error arised getting response for %v (%v)",
			request.URL.String(),
			err.Error())
		return nil, err
	}
	return response, nil
}
