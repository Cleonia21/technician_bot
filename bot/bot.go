package bot

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/withmandala/go-log"
	"main/conf"
	"main/db"
	"main/utils"
	"os"
)

type Bot struct {
	telegram *telego.Bot
	db       *db.Data
	logger   *log.Logger
}

func Init(dataBase *db.Data) *Bot {
	b := &Bot{} //
	botToken := conf.TOKEN

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	err := errors.New("")
	b.telegram, err = telego.NewBot(botToken) //, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.logger, _ = utils.NewLogger("")

	b.db = dataBase

	return b
}

func (b *Bot) Stop() {
	b.db.Close()
}

func (b *Bot) Start() {
	// Get tgGetChan channel
	tgGetChan, _ := b.telegram.UpdatesViaLongPolling(nil)

	// Stop reviving tgGetChan from update channel
	defer b.telegram.StopLongPolling()

	// Loop through all tgGetChan when they came
	for update := range tgGetChan {
		if update.CallbackQuery != nil {
			b.callbackQueryHandler(update.CallbackQuery)
		}
		if update.Message != nil && update.Message.Chat.Type == "private" {
			b.msgHandler(update.Message)
		}
	}
}
