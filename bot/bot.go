package bot

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"main/conf"
	"main/db"
	"os"
)

type Bot struct {
	telegram  *telego.Bot
	db        db.Data
	hashCache map[string]string
}

func Init() *Bot {
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

	b.db = db.Init(b.telegram.Logger())
	b.hashCache = make(map[string]string)

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
		if update.Message != nil {
			b.message(update.Message)
		}
		if update.CallbackQuery != nil {
			b.callbackQuery(update.CallbackQuery)
		}
	}
}
