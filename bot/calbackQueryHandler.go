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
}

func (b *Bot) queryDataToMenu(data string) (keyboard *telego.InlineKeyboardMarkup, text string, err error) {
	if data == "start" {
		keyboard, text = b.startMsgParams()
		return
	}

	table, targetKey, sourceKey, err := b.parseData(data)
	if err != nil {
		return
	}

	text, err = b.db.GetValue(table, targetKey)
	if err != nil {
		return
	}

	keyboard = b.getKeyboard(table, targetKey)
	keyboard = b.addControlBtns(keyboard, table, sourceKey)

	return
}

func (b *Bot) parseData(data string) (table, targetKey, sourceKey string, err error) {
	splitData := strings.Split(data, "@")
	if len(splitData) != 2 {
		err = errors.New("invalid query data")
		return
	}

	table = splitData[0]
	sourceKey = splitData[1]
	targetKey, err = b.db.GetTarget(table, sourceKey)
	if err != nil {
		err = fmt.Errorf("not found target %v", err)
		return
	}
	return
}

func (b *Bot) addControlBtns(keyboard *telego.InlineKeyboardMarkup, table string, sourceKey string) *telego.InlineKeyboardMarkup {
	if sourceKey == "start" {
		keyboard = addControlBtns(keyboard, "")
	} else {
		parent, err := b.db.GetParent(table, sourceKey)
		if err != nil {
			keyboard = addControlBtns(keyboard, "start")
			return keyboard
		}
		source, err := b.db.GetSource(table, parent)
		if err != nil {
			keyboard = addControlBtns(keyboard, "start")
			return keyboard
		}
		keyboard = addControlBtns(keyboard, fmt.Sprintf("%v@%v", table, source))
	}
	return keyboard
}

func addControlBtns(markup *telego.InlineKeyboardMarkup, key string) *telego.InlineKeyboardMarkup {
	startBtn := telegoutil.InlineKeyboardButton("в начало").WithCallbackData("start")
	if key == "" {
		markup.InlineKeyboard = append(markup.InlineKeyboard, telegoutil.InlineKeyboardRow(startBtn))
	} else {
		backBtn := telegoutil.InlineKeyboardButton("назад").WithCallbackData(key)
		markup.InlineKeyboard = append(markup.InlineKeyboard, telegoutil.InlineKeyboardRow(startBtn, backBtn))
	}
	return markup
}

func (b *Bot) getKeyboard(table string, key string) (keyboard *telego.InlineKeyboardMarkup) {
	btns, err := b.getBtns(table, key)
	if err != nil {
		return new(telego.InlineKeyboardMarkup)
	}

	params := b.getParams(table, key)

	keyboard, err = placementBtns(btns, params)
	if err != nil {
		b.logger.Error(err)
		return
	}
	return
}

func (b *Bot) getBtns(table string, key string) ([]telego.InlineKeyboardButton, error) {
	child, err := b.db.GetChild(table, key)
	if err != nil || len(child) == 0 {
		return nil, errors.New("buttons not found")
	}

	var btns []telego.InlineKeyboardButton
	for key, value := range child {
		if len(key) == 0 || len(value) == 0 {
			continue
		}
		btns = append(btns,
			telegoutil.InlineKeyboardButton(value).
				WithCallbackData(fmt.Sprintf("%v@%v", table, key)))
	}
	return btns, nil
}

func (b *Bot) getParams(table string, key string) string {
	var params string
	target, err := b.db.GetTarget(table, key)
	if err == nil {
		value, err := b.db.GetValue(table, target)
		if err == nil {
			params = value
		}
	}
	return params
}
