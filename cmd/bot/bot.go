package bot

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/withmandala/go-log"
	"os"
	"technician_bot/cmd/conf"
	"technician_bot/cmd/utils"
)

type Bot struct {
	telegram *telego.Bot
	logger   *log.Logger
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

	b.logger, _ = utils.NewLogger("")

	return b
}

func (b *Bot) Stop() {

}

func (b *Bot) Start() {
	// Get tgGetChan channel
	tgGetChan, _ := b.telegram.UpdatesViaLongPolling(nil)

	// Stop reviving tgGetChan from update channel
	defer b.telegram.StopLongPolling()

	for range [10]int{} {
		go func() {
			// Loop through all tgGetChan when they came
			for update := range tgGetChan {
				if update.CallbackQuery != nil {
					b.callbackQueryHandler(update.CallbackQuery)
				}
				if update.Message != nil && update.Message.Chat.Type == "private" {
					b.msgHandler(update.Message)
				}
			}
		}()
	}
	plug := make(chan int)
	<-plug
}
