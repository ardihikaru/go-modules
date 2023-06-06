package telegrambot

import (
	"fmt"
	"strings"
	"time"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/satumedishub/go-modules/pkg/enums/emoji"

	"github.com/satumedishub/go-modules/pkg/enums/loglevel"
	"github.com/satumedishub/go-modules/pkg/logger"
	m "github.com/satumedishub/go-modules/pkg/messenger"
)

type TelegramBot struct {
	Messenger   *m.Messenger
	Bot         *tgBotApi.BotAPI
	updates     tgBotApi.UpdatesChannel
	log         *logger.Logger
	groupChatId int64
	groupTitle  string
}

// Connect initializes telegram bot via API
func Connect(log *logger.Logger, token, logLevel string, groupTitle *string, groupChatId *int64,
	messenger *m.Messenger) (*TelegramBot, error) {
	bot, err := tgBotApi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// sets bot log level
	if loglevel.DebugMode(logLevel) {
		bot.Debug = true
	}

	log.Info(fmt.Sprintf("Authorized on Telegram account [%s]", bot.Self.UserName))

	if groupTitle == nil && groupChatId != nil && messenger != nil {
		// only groupChatId and messenger exists
		return &TelegramBot{Bot: bot, log: log, groupTitle: "", groupChatId: *groupChatId, Messenger: messenger}, nil
	} else if groupTitle == nil && groupChatId != nil && messenger == nil {
		// only groupChatId exists
		return &TelegramBot{Bot: bot, log: log, groupTitle: "", groupChatId: *groupChatId}, nil
	} else if groupTitle != nil && groupChatId != nil && messenger != nil {
		// all exists
		return &TelegramBot{Bot: bot, log: log, groupTitle: *groupTitle, groupChatId: *groupChatId,
			Messenger: messenger}, nil
	} else {
		return &TelegramBot{Bot: bot, log: log}, nil
	}
}

// Init initializes telegram messages
func (b *TelegramBot) Init() {
	b.log.Info(fmt.Sprintf("Authorized on account [%s]", b.Bot.Self.UserName))

	u := tgBotApi.NewUpdate(0)
	u.Timeout = 60

	b.updates = b.Bot.GetUpdatesChan(u)
}

// Run starts subscribing messages
func (b *TelegramBot) Run() {
	for update := range b.updates {
		if update.Message != nil {
			var err error

			// only response message that has originated message
			replyToMsg := update.Message.ReplyToMessage
			if replyToMsg == nil {
				continue
			}

			// only response message that's originated from the Bot
			repliedMsgOwnerAsBot := replyToMsg.From.IsBot
			if !repliedMsgOwnerAsBot {
				continue
			}
			repliedMsg := cleanMsg(replyToMsg.Text)

			// at least 5 slices is a considered as a valid message
			repliedMsgList := strings.Split(repliedMsg, "|")
			if len(repliedMsgList) < 5 {
				continue
			}

			// extracts important information
			phone := repliedMsgList[2]

			// identifies the source of message
			msgType := update.Message.Chat.Type    // e.g. "group"
			groupName := update.Message.Chat.Title // e.g. "Chats CS SatuMedis PROD"
			if msgType != "group" && groupName != b.groupTitle {
				b.log.Debug("found ignorable messages. ignore captured message")
				continue
			}

			// process sending message to the Messenger
			sent, msgToReply, err := b.Messenger.SendMsgToWhatsapp(phone, update.Message.Text)
			if !sent {
				msgToReply = fmt.Sprintf("failed to reply chat from recipient [%s]", phone)
				return
			}

			msg := tgBotApi.NewMessage(update.Message.Chat.ID, msgToReply)
			msg.ReplyToMessageID = update.Message.MessageID

			err = send(b.Bot, msg)
			if err != nil {
				b.log.Warn(fmt.Sprintf("failed to send the message -> %s", err.Error()))
				return
			}
		}
	}
}

// cleanMsg removes unexpected strings
func cleanMsg(msg string) string {
	msg = strings.Replace(msg, "SasaBot|", "", len(msg))
	msg = strings.Replace(msg, "MESSAGE_ID:", "", len(msg))
	msg = strings.Replace(msg, "PHONE:", "", len(msg))
	msg = strings.Replace(msg, "FROM:", "", len(msg))
	msg = strings.Replace(msg, fmt.Sprintf(" %s ", emoji.CheckMark), "", len(msg))
	msg = strings.Replace(msg, "MESSAGE:", "", len(msg))
	msg = strings.Replace(msg, "\n", "", len(msg))

	return msg
}

// send sends a text message
func send(bot *tgBotApi.BotAPI, msg tgBotApi.MessageConfig) error {
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// SendTextMsg sends messages to the telegram bot
func (b *TelegramBot) SendTextMsg(ts time.Time, phone, msgId, name, msg string) {
	// builds telegram message content
	msgTemplate := buildTelegramMessage(ts, msgId, phone, name, msg)

	// sends the message
	tgMessage := tgBotApi.NewMessage(b.groupChatId, msgTemplate)
	_, err := b.Bot.Send(tgMessage)
	if err != nil {
		b.log.Error(fmt.Sprintf("sending Telegram message failed -> %s", err.Error()))
		return
	}
}

// buildTelegramMessage generates telegram formatted message with a predefined message template
func buildTelegramMessage(ts time.Time, msgId, phone, name, msg string) string {
	return fmt.Sprintf("SasaBot|%s|\n"+
		"MESSAGE_ID:%s|\n"+
		"PHONE:%s|\n"+
		"FROM:%s %s |\n\n"+
		"MESSAGE:\n%s",
		ts.String(), msgId, phone, name, emoji.CheckMark, msg)
}
