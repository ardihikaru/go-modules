package whatsappbot

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	"github.com/satumedishub/go-modules/pkg/logger"
	tgBot "github.com/satumedishub/go-modules/pkg/telegrambot"
)

type DataPasser struct {
	logs chan string
}

type WhatsappBot struct {
	TelegramBot    *tgBot.TelegramBot
	Client         *whatsmeow.Client
	log            *logger.Logger
	eventHandlerID uint32
}

// Connect builds whatsapp client and connects to it
func Connect(telegramBot *tgBot.TelegramBot, dbName string, log *logger.Logger) (*WhatsappBot, error) {
	address := fmt.Sprintf("file:%s.db?_foreign_keys=on", dbName)

	container, err := sqlstore.New("sqlite3", address, waLog.Noop)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	client := whatsmeow.NewClient(deviceStore, waLog.Noop)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				log.Info(fmt.Sprintf("Login event: %s", evt.Event))
			}
		}
	} else {
		err := client.Connect()
		if err != nil {
			return nil, err
		}

		log.Info("whatsapp account has been connected successfully")
	}
	return &WhatsappBot{TelegramBot: telegramBot, Client: client, log: log}, nil
}

// SendMsg sends message to designated whatsapp number
func (wb *WhatsappBot) SendMsg(recipient types.JID, phone, msg string) error {
	resp, err := wb.Client.SendMessage(context.Background(), recipient, &waProto.Message{
		Conversation: proto.String(msg),
	})
	if err != nil {
		wb.log.Debug(fmt.Sprintf("[to:%s] failed to send message: %s", phone, msg))
		return err
	} else {
		wb.log.Debug(fmt.Sprintf("[to:%s] message sent (server timestamp: %s)", phone, resp.Timestamp))
	}

	return nil
}

// Register registers a new event handler
func (wb *WhatsappBot) Register() {
	wb.eventHandlerID = wb.Client.AddEventHandler(wb.eventHandler)
}

// eventHandler handles incoming events from the whatsapp chat
func (wb *WhatsappBot) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		msgId := v.Info.ID
		phone := v.Info.Sender.User
		name := v.Info.PushName
		ts := v.Info.Timestamp
		message := v.Message.GetConversation()

		// monkey patch! somethimes the text is not in the Conversation, but in the ExtendedTextMessage
		// e.g. from Albert / Taiwan
		if message == "" {
			message = *v.Message.ExtendedTextMessage.Text
		}

		wb.log.Debug(fmt.Sprintf("**** [%s][%s] Received a message from [%s] (%s)! -> '%s'\n\n",
			ts, msgId, name, phone, message))

		// sends to telegram messenger
		wb.TelegramBot.SendTextMsg(ts, phone, msgId, name, message)
	}
}
