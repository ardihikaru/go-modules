package wawebhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/whatsmeow/types"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
)

type ReplyMessage struct {
	Message       string
	WithImage     bool
	ImageFileName string
}

// replyMessage replies the captured message and do reply
func (wb *WaBot) replyMessage(targetJID *types.JID, phone string, resp *httputils.Response) error {
	// extracts response payload
	byteData, _ := json.Marshal(resp.Data)
	replyMsgObj := ReplyMessage{}
	err := json.Unmarshal(byteData, &replyMsgObj)
	if err != nil {
		wb.Log.Error("failed to convert reply payload response", zap.Error(err))
		return err
	}

	if targetJID.User == "" {
		// enriches with `+` symbol if missing
		if phone[0:1] != "+" {
			phone = fmt.Sprintf("+%s", phone)
		}

		// validates phone number and get the recipient
		recipient, err := wb.ValidateAndGetRecipient(phone, true)
		if err != nil {
			wb.Log.Error(fmt.Sprintf("phone [%s] got validation error(s)", phone), zap.Error(err))
			return err
		}

		// sends a reply message
		err = wb.sendMsgAndWait(*recipient, replyMsgObj)
		if err != nil {
			wb.Log.Error("failed to reply the captured message", zap.Error(err))
			return err
		}

	} else {
		// TODO: process outgoing message
	}

	return nil
}

// ValidateAndGetRecipient validates the phone number
func (wb *WaBot) ValidateAndGetRecipient(phone string, ignoreInContactList bool) (*types.JID, error) {
	var err error
	phones := make([]string, 1)

	// enriches with `+` symbol if missing
	if phone[0:1] != "+" {
		phone = fmt.Sprintf("+%s", phone)
	}

	phones[0] = phone

	// checks if this number available on Whatsapp or not
	onWhatsapp, err := wb.Client.IsOnWhatsApp(phones)
	if err != nil {
		return nil, err
	}
	wb.Log.Debug(fmt.Sprintf("%v", onWhatsapp))

	// extracts non-AD JID
	recipient := onWhatsapp[0].JID
	if !onWhatsapp[0].IsIn {
		wb.Log.Warn(fmt.Sprintf("this number [%s] is not available in Whatsapp", phone))
		return nil, fmt.Errorf("this number is not available in Whatsapp")
	}

	// check if in contact list
	// TODO: if not in the contact list, what to do?
	contact, err := wb.Client.Store.Contacts.GetContact(recipient)
	if err != nil && !ignoreInContactList {
		wb.Log.Warn(fmt.Sprintf("this number [%s] is not on your contact list!", phone))
		return nil, err
	}
	wb.Log.Debug(fmt.Sprintf("%v", contact))

	return &recipient, nil
}

// sendMsgAndWait sends the message to the designated device
func (wb *WaBot) sendMsgAndWait(recipient types.JID, msgObj ReplyMessage) error {
	var err error

	if msgObj.WithImage {
		imgPath := fmt.Sprintf("%s/%s", wb.ImageDir, msgObj.ImageFileName)
		imgInBytes, uploaded, err := wb.UploadImgToWhatsapp(imgPath)

		// prepares image information
		contentType := http.DetectContentType(*imgInBytes)
		fileLength := uint64(len(*imgInBytes))

		if err != nil {
			return err
		}
		err = wb.SendImgMsg(recipient, uploaded, msgObj.Message, contentType, fileLength)
	} else {
		err = wb.SendMsg(recipient, msgObj.Message)
	}
	if err != nil {
		return err
	}

	return nil
}
