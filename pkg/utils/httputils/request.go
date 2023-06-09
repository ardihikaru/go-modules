package httputils

import (
	"encoding/json"
	"io"
	"net/http"
)

// GetJsonBody extracts json data from the request payload
func GetJsonBody(rBody io.ReadCloser, destType interface{}) (int, int, error) {
	// extracts request body
	b, err := io.ReadAll(rBody)
	if err != nil {
		return InvalidRequestJSON, http.StatusBadRequest, err
	}

	// read JSON body from the request
	err = json.Unmarshal(b, &destType)
	if err != nil {
		return RequestJSONExtractionFailed, http.StatusBadRequest, err
	}

	return 0, 200, nil
}
