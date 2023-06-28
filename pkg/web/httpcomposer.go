package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
)

// BuildRequest builds new request
func BuildRequest(apiUrl, method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, apiUrl, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ExecRequest executes request
func ExecRequest(httpClient *http.Client, req *http.Request) (interface{}, error) {
	// sends request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// checks status code, if not 200, returns the error message from the RESTApi
	// otherwise, it means that the JSON data extraction was failed
	if resp.StatusCode == 204 {
		return nil, fmt.Errorf("no content found")
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("designated API not found")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch data from the designated service")
	}

	// read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(httputils.ResponseText("", httputils.InvalidRequestJSON)), err
	}

	// Convert response body to designated struct (respPayload)
	var respPayload interface{}
	err = json.Unmarshal(bodyBytes, &respPayload)
	if err != nil {
		return nil, fmt.Errorf(httputils.ResponseText("", httputils.RequestJSONExtractionFailed))
	}

	return respPayload, nil
}
