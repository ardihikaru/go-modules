package qrcodehandler

import (
	"github.com/yougg/go-qrcode"
)

// StoreQrCode stores QR Code in a file and send to the target
func StoreQrCode(code, filePath string) error {
	err := qrcode.WriteFile(code, qrcode.Medium, 256, filePath, 0)
	if err != nil {
		return err
	}

	return nil
}
