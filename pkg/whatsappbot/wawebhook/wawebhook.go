package wawebhook

import (
	"context"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-modules/pkg/logger"
	fh "github.com/ardihikaru/go-modules/pkg/utils/filehandler"
	qrCodeH "github.com/ardihikaru/go-modules/pkg/utils/qrcodehandler"
)

const (
	IncomingMessage = "INCOMING_MESSAGE"
	OutgoingMessage = "OUTGOING_MESSAGE"
)

// WaManager defines the whatsapp container
type WaManager struct {
	Container *sqlstore.Container
	Log       *logger.Logger
}

// WaBot defines the bot client
type WaBot struct {
	Client         *whatsmeow.Client
	Log            *logger.Logger
	HttpClient     *http.Client
	EventHandlerID uint32
	Phone          string
	WebhookUrl     string
	ImageDir       string
	EchoMsg        bool
	WHookEnabled   bool
}

// BotClientList defines the variable to store WaBot objects
type BotClientList map[string]*WaBot

// NewContainer builds whatsapp Container
func NewContainer(dbName string, log *logger.Logger) (*WaManager, error) {
	address := fmt.Sprintf("file:%s.db?_foreign_keys=on", dbName)

	container, err := sqlstore.New("sqlite3", address, waLog.Noop)
	if err != nil {
		return nil, err
	}

	return &WaManager{Container: container, Log: log}, nil
}

// LoginExistingWASession logins with an existing session on the database
func LoginExistingWASession(httpClient *http.Client, webhookUrl, imageDir string, container *sqlstore.Container,
	log *logger.Logger, jidStr, phone string, echoMsg, wHookEnabled bool) (*WaBot, error) {
	// build JID
	jid, _ := types.ParseJID(jidStr)

	deviceStore, _ := container.GetDevice(jid)
	if deviceStore == nil {
		return nil, fmt.Errorf("unable to find device with this JID")
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	err := client.Connect()
	if err != nil {
		log.Error("WhatsApp conn error", zap.Error(err))
		return nil, fmt.Errorf("unable connect")
	}

	return buildWhatsappBot(client, log, httpClient, phone, webhookUrl, imageDir, echoMsg, wHookEnabled), nil
}

// NewWhatsappClient initializes Whatsapp client
func NewWhatsappClient(httpClient *http.Client, webhookUrl, imageDir string, container *sqlstore.Container, log *logger.Logger,
	phone, fileDir string, echoMsg, wHookEnabled bool, printTerminal bool) (*WaBot, error) {
	var err error

	myDevice := container.NewDevice()
	clientLog := waLog.Stdout("Client", "INFO", true)

	// generates a new client
	client := whatsmeow.NewClient(myDevice, clientLog)

	// generates file path to store the qr code
	filePath := fmt.Sprintf("%s/%s.png", fileDir, phone)

	// No ID stored, new login
	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		return nil, err
	}

	// sending qrCode
	for evt := range qrChan {
		if evt.Event == "code" {
			err = qrCodeH.StoreQrCode(evt.Code, filePath)
			if err != nil {
				log.Error("failed to write QR Code to file", zap.Error(err))
				return nil, err
			}

			// prints qrcode in terminal (if enabled)
			if printTerminal {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			}
		} else {
			log.Info(fmt.Sprintf("Login event: %s", evt.Event))

			// since there is no valid open client, cancel the request
			// return nil, fmt.Errorf("failed to open a new session")
		}
	}

	// new session has been created. deleting the file
	// ignores error event if happens (e.g. ignore if file does not exists)
	_ = fh.DeleteFile(filePath)

	return buildWhatsappBot(client, log, httpClient, phone, webhookUrl, imageDir, echoMsg, wHookEnabled), nil
}

func buildWhatsappBot(client *whatsmeow.Client, log *logger.Logger, httpClient *http.Client,
	phone, webhookUrl, imageDir string, echoMsg, wHookEnabled bool) *WaBot {
	return &WaBot{
		Client:       client,
		Log:          log,
		Phone:        phone,
		HttpClient:   httpClient,
		WebhookUrl:   webhookUrl,
		ImageDir:     imageDir,
		EchoMsg:      echoMsg,
		WHookEnabled: wHookEnabled,
	}
}

// Register registers a new event handler
func (wb *WaBot) Register() {
	wb.EventHandlerID = wb.Client.AddEventHandler(wb.eventHandler)
}

// eventHandler handles incoming events from the whatsapp chat
func (wb *WaBot) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// when DeviceSentMeta is nil -> Incoming message (this device is receiving a message)
		// when DeviceSentMeta is NOT nil -> Outgoing message (this device is sending a message)
		var deviceTargetJID types.JID
		if v.Info.DeviceSentMeta != nil {
			deviceTargetJIDStr := v.Info.DeviceSentMeta.DestinationJID
			deviceTargetJID, _ = types.ParseJID(deviceTargetJIDStr)
		}

		msgId := v.Info.ID
		msgType := v.Info.Type // e.g. Text
		phone := v.Info.Sender.User
		name := v.Info.PushName
		ts := v.Info.Timestamp
		message := v.Message.GetConversation()

		// do nothing if webhook disabled
		if !wb.WHookEnabled {
			return
		}

		// do nothing if webhook URL is empty
		if wb.WebhookUrl == "" {
			wb.Log.Warn("invalid Webhook URL due to an empty value")
			return
		}

		// monkey patch! sometimes the text is not in the Conversation, but in the ExtendedTextMessage
		// e.g. from Albert / Taiwan
		if message == "" && v.Message.ExtendedTextMessage != nil {
			message = *v.Message.ExtendedTextMessage.Text
		}

		if message != "" && v.Info.DeviceSentMeta == nil {
			wb.Log.Info(fmt.Sprintf("**** [%s][%s] Received a [%s] message from [%s] (%s) -> '%s'",
				ts, msgId, msgType, name, phone, message))

			// on receiving message, send the message to the designated webhook
			resp, err := wb.sendToWebhook(&deviceTargetJID, IncomingMessage, msgId, msgType, phone, name, message,
				wb.Phone, ts)
			if err != nil {
				wb.Log.Error("failed to forward incoming message to webhook", zap.Error(err))
			} else {
				wb.replyMessage(&deviceTargetJID, phone, resp)
			}
		} else if message != "" && v.Info.DeviceSentMeta != nil {
			//} else if message != "" && phone == wb.Phone {
			wb.Log.Debug(fmt.Sprintf("**** [%s][%s] Sent a [%s] message to [%s] (%s) -> '%s'",
				ts, msgId, msgType, name, deviceTargetJID.User, message))

			// TODO: it sends a message to itself, any special action to do?
			if phone == wb.Phone {
			}

			// TODO: on sent message, do something here
			//_, err := wb.sendToWebhook(nil, OutgoingMessage, msgId, msgType, phone, name, message, ts)
			//if err != nil {
			//	wb.Log.Error("failed to forward outgoing message to webhook", zap.Error(err))
			//}
		}
	}
}
