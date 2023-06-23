package wawebhook

import "github.com/ardihikaru/go-modules/pkg/utils/common"

type MessagePayload struct {
	From          string `json:"from"`
	To            string `json:"to"`
	Message       string `json:"message,omitempty"`
	ImageFileName string `json:"image_filename,omitempty"`
	ImageCaption  string `json:"image_caption,omitempty"`
}

// Validate validates message payload
func (p *MessagePayload) Validate() error {
	return nil
}

// Sanitize sanitizes message payload
func (p *MessagePayload) Sanitize() {
	plusSymbol := false

	p.From = common.SanitizePhone(p.From, &plusSymbol)
	p.To = common.SanitizePhone(p.To, &plusSymbol)
}
