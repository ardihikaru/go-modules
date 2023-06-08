package whatsappbot

import (
	"github.com/ardihikaru/go-modules/pkg/logger"
	e "github.com/ardihikaru/go-modules/pkg/utils/error"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"
)

// InitWhatsappContainer initializes Whatsapp Container
func InitWhatsappContainer(dbName string, log *logger.Logger) *botHook.WaManager {
	waManager, err := botHook.NewContainer(dbName, log)
	if err != nil {
		e.FatalOnError(err, "failed to initialize WhatsApp Bot Container")
	}

	return waManager
}
