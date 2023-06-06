package messenger

import (
	"crypto/tls"
	"net/http"
)

// BuildHttpClient builds http client
func BuildHttpClient(httpClientTLS bool) *http.Client {
	if httpClientTLS {
		// enriches with TLS information if enabled
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}

	} else {
		return &http.Client{}
	}
}
