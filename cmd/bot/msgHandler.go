package bot

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"technician_bot/database"
)

func (b *Bot) msgHandler(msg *telego.Message) {
	var keyboard *telego.InlineKeyboardMarkup
	var text string

	if msg.Text == "/start" {
		keyboard, text = b.startMsgParams()
	} else {
		return
	}
	b.sendMenu(telegoutil.ID(msg.Chat.ID), "", keyboard, text)
}

func (b *Bot) startMsgParams() (keyboard *telego.InlineKeyboardMarkup, text string) {
	tables := b.tableFirsKeys()
	var btns []telego.InlineKeyboardButton

	for key, value := range tables {
		btns = append(btns, telegoutil.InlineKeyboardButton(key).
			WithCallbackData(value))
	}
	keyboard = btnsOptimalPlacement(btns)
	text = "Привет! Выбери машину"
	return keyboard, text
}

func (b *Bot) tableFirsKeys() map[string]string {
	tables := make(map[string]string)

	keyKia, err := database.GetKey("kia", "start")
	if err != nil {
		b.logger.Error("table or first row kia not found")
	}
	tables["kia"] = "kia@" + keyKia

	keyPolo, err := database.GetKey("polo", "start")
	if err != nil {
		b.logger.Error("table or first row polo not found")
	}
	tables["polo"] = "polo@" + keyPolo

	return tables
}
