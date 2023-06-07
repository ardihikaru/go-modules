package wawebhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	fh "github.com/ardihikaru/go-modules/pkg/utils/filehandler"
	qrCodeH "github.com/ardihikaru/go-modules/pkg/utils/qrcodehandler"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"
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
	EventHandlerID uint32
	Phone          string
	WebhookUrl     string
	HttpClient     *http.Client
}

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
func LoginExistingWASession(httpClient *http.Client, webhookUrl string, container *sqlstore.Container, log *logger.Logger, jidStr, phone string) (*WaBot, error) {
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

	return buildWhatsappBot(client, log, phone, httpClient, webhookUrl), nil
}

// NewWhatsappClient initializes Whatsapp client
func NewWhatsappClient(httpClient *http.Client, webhookUrl string, container *sqlstore.Container, log *logger.Logger,
	phone, fileDir string) (*WaBot, error) {
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
		} else {
			log.Info(fmt.Sprintf("Login event: %s", evt.Event))
		}
	}

	// new session has been created. deleting the file
	// ignores error event if happens (e.g. ignore if file does not exists)
	_ = fh.DeleteFile(filePath)

	return buildWhatsappBot(client, log, phone, httpClient, webhookUrl), nil
}

func buildWhatsappBot(client *whatsmeow.Client, log *logger.Logger, phone string, httpClient *http.Client,
	webhookUrl string) *WaBot {
	return &WaBot{
		Client:     client,
		Log:        log,
		Phone:      phone,
		HttpClient: httpClient,
		WebhookUrl: webhookUrl,
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
		msgId := v.Info.ID
		phone := v.Info.Sender.User
		name := v.Info.PushName
		ts := v.Info.Timestamp
		message := v.Message.GetConversation()

		// monkey patch! sometimes the text is not in the Conversation, but in the ExtendedTextMessage
		// e.g. from Albert / Taiwan
		if message == "" && v.Message.ExtendedTextMessage != nil {
			message = *v.Message.ExtendedTextMessage.Text
		}

		if message != "" {
			wb.Log.Info(fmt.Sprintf("**** [%s][%s] Received a message from [%s] (%s)! -> '%s'",
				ts, msgId, name, phone, message))

			// on receiving message, send the message to the designated webhook
		}
	}
}
