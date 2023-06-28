package wawebhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ardihikaru/go-modules/pkg/utils/web"
	"io"
	"os"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
)

type WebhookBody struct {
	PhoneOwner   string `json:"phone_owner"`
	EventType    string `json:"event_type"`
	MsgId        string `json:"msg_id"`
	MsgType      string `json:"msg_type"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	TargetJID    string `json:"target_jid"`
	TargetDevice string `json:"target_device"`
	Timestamp    string `json:"timestamp"`
}

// sendToWebhook sends the captured message to webhook
func (wb *WaBot) sendToWebhook(targetJID *types.JID, evtType, msgId, msgType, phone, name, message, phoneOwner string,
	ts time.Time) (*httputils.Response, error) {
	// if echo message enabled, simply send and echo message
	if wb.EchoMsg {
		return &httputils.Response{
			Data: ReplyMessage{
				Message: message,
			},
		}, nil
	} else {
		// otherwise, send POST request to the designated webhook
		// prepare body
		bodyObj := &WebhookBody{
			PhoneOwner:   phoneOwner,
			EventType:    evtType,
			MsgId:        msgId,
			MsgType:      msgType,
			Phone:        phone,
			Name:         name,
			Message:      message,
			TargetJID:    targetJID.String(),
			TargetDevice: targetJID.User,
			Timestamp:    ts.Format("2006-01-02 15:04:05"),
		}

		// builds body
		body, err := web.BuildFormBody(bodyObj)

		// builds request
		req, err := web.BuildRequest(wb.WebhookUrl, "POST", body)
		if err != nil {
			return nil, err
		}

		// enriches with pre-generated headers
		req.Header.Set("Content-Type", "application/json") // default header

		// sends request
		resp, err := wb.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// validates response
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("got error response from the webhook")
		}

		// read response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// converts response body to clinicRespPayload struct
		var respPayload httputils.Response
		err = json.Unmarshal(bodyBytes, &respPayload)
		if err != nil {
			return nil, err
		}

		return &respPayload, nil
	}
}

// SendMsg sends message to designated whatsapp number
func (wb *WaBot) SendMsg(recipient types.JID, msg string) error {
	resp, err := wb.Client.SendMessage(context.Background(), recipient, &waProto.Message{
		Conversation: proto.String(msg),
	})
	if err != nil {
		wb.Log.Debug(fmt.Sprintf("[to:%s] failed to send message: %s", recipient.User, msg))
		return err
	} else {
		wb.Log.Debug(fmt.Sprintf("[to:%s] message sent (server timestamp: %s)", recipient.User, resp.Timestamp))
	}

	return nil
}

// SendImgMsg sends image-based message to designated whatsapp number
func (wb *WaBot) SendImgMsg(recipient types.JID, uploadedImg *whatsmeow.UploadResponse, imgCaption, contentType string,
	fileLength uint64) error {
	msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
		Caption:       proto.String(imgCaption),
		Url:           proto.String(uploadedImg.URL),
		DirectPath:    proto.String(uploadedImg.DirectPath),
		MediaKey:      uploadedImg.MediaKey,
		Mimetype:      proto.String(contentType),
		FileEncSha256: uploadedImg.FileEncSHA256,
		FileSha256:    uploadedImg.FileSHA256,
		FileLength:    proto.Uint64(fileLength),
	}}

	resp, err := wb.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		wb.Log.Debug(fmt.Sprintf("[to:%s] failed to send image message: %s", recipient.User, msg))
		return err
	} else {
		wb.Log.Debug(fmt.Sprintf("[to:%s] message sent (server timestamp: %s)", recipient.User, resp.Timestamp))
	}

	return nil
}

// UploadImgToWhatsapp uploads the prepared image to Whatsapp server
func (wb *WaBot) UploadImgToWhatsapp(imgPath string) (*[]byte, *whatsmeow.UploadResponse, error) {
	// first, prepares the image file as bytes
	imgInBytes, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, nil, err
	}

	// then, uploads to whatsapp server
	uploaded, err := wb.Client.Upload(context.Background(), imgInBytes, whatsmeow.MediaImage)
	if err != nil {
		wb.Log.Debug(fmt.Sprintf("[to:%s] failed to upload file"))
		return nil, nil, err
	}

	return &imgInBytes, &uploaded, nil
}
