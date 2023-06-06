package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/satumedishub/go-modules/pkg/logger"
	"github.com/satumedishub/go-modules/pkg/web"
)

type Messenger struct {
	Log        *logger.Logger
	HttpClient *http.Client
	Url        string
}

// PostWhatsappMsg defines the field parameters
type PostWhatsappMsg struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// WhatsappResponse defines the field parameters
type WhatsappResponse struct {
	Code     int64       `json:"code,omitempty"`
	ErrorMsg string      `json:"error,omitempty"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

// InitMessenger instantiates messenger
func InitMessenger(url string, tls bool, log *logger.Logger) *Messenger {
	return &Messenger{
		HttpClient: BuildHttpClient(tls),
		Log:        log,
		Url:        url,
	}
}

// SendMsgToWhatsapp sends messages to the Whatsapp chat
func (m *Messenger) SendMsgToWhatsapp(phone, msg string) (bool, string, error) {
	var err error
	var respPayload WhatsappResponse

	// builds auth post form
	payload := PostWhatsappMsg{
		Phone:   phone,
		Message: msg,
	}

	// prepares the POST form in bytes
	authFormBytes := new(bytes.Buffer)
	err = json.NewEncoder(authFormBytes).Encode(payload)
	if err != nil {
		return false, "", err
	}

	// builds POST request with JSON payload
	req, err := http.NewRequest("POST", m.Url, authFormBytes)
	if err != nil {
		return false, "", err
	}
	req.Header.Set(web.HeaderContentTypeKey, web.HeaderContentTypeValue)

	// sends request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()

	// read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	// Convert response body to WhatsappResponse struct
	err = json.Unmarshal(bodyBytes, &respPayload)
	if err != nil {
		return false, "", err
	}

	// if data is nil, returns an error
	if respPayload.Data == nil {
		return false, "", fmt.Errorf("response data is empty")
	}

	return true, "chat has been replied", nil
}
