package wawebhook

import (
	"encoding/json"
	"fmt"

	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	"go.mau.fi/whatsmeow/types"
	"go.uber.org/zap"
)

type ReplyMessage struct {
	Message string
}

// replyMessage replies the captured message and do reply
func (wb *WaBot) replyMessage(targetJID *types.JID, phone string, resp *httputils.Response) {
	// extracts response payload
	byteData, _ := json.Marshal(resp.Data)
	replyMsgObj := ReplyMessage{}
	err := json.Unmarshal(byteData, &replyMsgObj)
	if err != nil {
		wb.Log.Error("failed to convert reply payload response", zap.Error(err))
		return
	}

	if targetJID.User == "" {
		// enriches with `+` symbol if missing
		if phone[0:1] != "+" {
			phone = fmt.Sprintf("+%s", phone)
		}

		// validates phone number
		recipient, err := wb.validateAndGetRecipient(phone, true)
		if err != nil {
			wb.Log.Error(fmt.Sprintf("phone [%s] got validation error(s)", phone), zap.Error(err))
			return
		}

		// sends a reply message
		err = wb.sendMsgAndWait(*recipient, replyMsgObj.Message)
		if err != nil {
			wb.Log.Error("failed to reply the captured message", zap.Error(err))
			return
		}

	} else {
		// TODO: process outgoing message
	}
}

// validate validates the phone number
func (wb *WaBot) validateAndGetRecipient(phone string, ignoreInContactList bool) (*types.JID, error) {
	var err error
	phones := make([]string, 1)

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
func (wb *WaBot) sendMsgAndWait(recipient types.JID, msg string) error {
	var err error

	err = wb.SendMsg(recipient, msg)
	if err != nil {
		return err
	}

	return nil
}
