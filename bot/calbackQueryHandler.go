package bot

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"strings"
)

func (b *Bot) callbackQueryHandler(query *telego.CallbackQuery) {
	keyboard, text, err := b.queryDataToMenu(query.Data)

	if err != nil {
		b.telegram.Logger().Errorf(err.Error())
		return
	}
	b.sendMenu(telegoutil.ID(query.Message.Chat.ID), query.ID, keyboard, text)
	//b.editMenu(query, keyboard, text)
}

func (b *Bot) queryDataToMenu(data string) (keyboard *telego.InlineKeyboardMarkup, text string, err error) {
	splitData := strings.Split(data, "@")
	if len(splitData) != 2 {
		return nil, "", errors.New("invalid query data")
	}

	table := splitData[0]
	sourceKey := splitData[1]
	targetKey, err := b.db.GetTarget(table, sourceKey)
	if err != nil {
		return nil, "", fmt.Errorf("not found target %v", err)
	}

	text, err = b.db.GetValue(table, targetKey)
	if err != nil {
		return
	}

	keyboard = b.keyToKeyboard(table, targetKey)

	return
}

func (b *Bot) keyToKeyboard(table string, parentKey string) (keyboard *telego.InlineKeyboardMarkup) {
	child, err := b.db.GetChild(table, parentKey)
	if err != nil {
		return nil
	}

	var btns []telego.InlineKeyboardButton
	for key, value := range child {
		btns = append(btns,
			telegoutil.InlineKeyboardButton(value).
				WithCallbackData(fmt.Sprintf("%v@%v", table, key)))
	}
	return b.btnsOptimalPlacement(btns)
}
