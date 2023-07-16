package bot

import (
	"errors"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"main/db"
	"strings"
)

func (b *Bot) message(msg *telego.Message) {
	var keyboard *telego.InlineKeyboardMarkup
	var text string
	if msg.Text == "/start" {
		keyboard, text = b.startMsgParams()
	} else {
		return
	}
	b.sendMenu(telegoutil.ID(msg.Chat.ID), "", keyboard, text)
}

func (b *Bot) sendMenu(chatID telego.ChatID, answerForQueryID string, keyboard *telego.InlineKeyboardMarkup, text string) {
	message := telegoutil.Message(
		chatID,
		text,
	).WithReplyMarkup(keyboard)
	_, err := b.telegram.SendMessage(message)
	if err != nil {
		b.telegram.Logger().Errorf(err.Error())
	}
	if answerForQueryID != "" {
		err = b.telegram.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{CallbackQueryID: answerForQueryID})
		if err != nil {
			b.telegram.Logger().Errorf(err.Error())
		}
	}

}

func (b *Bot) callbackQuery(query *telego.CallbackQuery) {
	var keyboard *telego.InlineKeyboardMarkup
	var text string
	var err error
	if b.hashToText(query.Data) == "/start" {
		keyboard, text = b.startMsgParams()
	} else {
		keyboard, text, err = b.msgParamsFromDB(query)
	}
	if err != nil {
		b.telegram.Logger().Errorf(err.Error())
		return
	}
	b.sendMenu(telegoutil.ID(query.Message.Chat.ID), query.ID, keyboard, text)
	//b.editMenu(query, keyboard, text)
}

func (b *Bot) startMsgParams() (keyboard *telego.InlineKeyboardMarkup, text string) {
	keyboard = telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("Kia x-line").
				WithCallbackData(b.textToHash("kiaxline;/"))),
	)
	text = "Привет! Выбери машину"
	return keyboard, text
}

func (b *Bot) editMenu(query *telego.CallbackQuery, keyboard *telego.InlineKeyboardMarkup, text string) {
	var msgParams telego.EditMessageTextParams
	msgParams.WithText(text)
	msgParams.WithMessageID(query.Message.MessageID)
	msgParams.WithChatID(telegoutil.ID(query.Message.Chat.ID))
	msgParams.WithReplyMarkup(keyboard)
	_, err := b.telegram.EditMessageText(&msgParams)
	if err != nil {
		b.telegram.Logger().Errorf(err.Error())
	}
}

func (b *Bot) msgParamsFromDB(query *telego.CallbackQuery) (keyboard *telego.InlineKeyboardMarkup, text string, err error) {
	queryData := strings.Split(b.hashToText(query.Data), ";")
	if len(queryData) != 2 {
		err = errors.New("not valid query data")
		return
	}
	ctn, err := b.getButtonData(queryData[0], queryData[1])
	if err != nil {
		return
	}
	var btns []telego.InlineKeyboardButton
	for _, childName := range ctn.Child {
		btns = append(btns,
			telegoutil.InlineKeyboardButton(childName).
				WithCallbackData(b.textToHash(b.hashToText(query.Data)+"/"+childName)))
	}
	keyboard = b.btnsOptimalPlacement(btns)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, telegoutil.InlineKeyboardRow(
		telegoutil.InlineKeyboardButton("вернуться в главное меню").
			WithCallbackData(b.textToHash("/start"))))
	text = ctn.Text
	return
}

// 50 символов в строке
// 2 символа между двумя кнопками
func (b *Bot) btnsOptimalPlacement(btns []telego.InlineKeyboardButton) *telego.InlineKeyboardMarkup {
	if len(btns) == 0 {
		return &telego.InlineKeyboardMarkup{}
	}
	maxTextLen := 50

	keyboard := &telego.InlineKeyboardMarkup{}
	var row []telego.InlineKeyboardButton
	var rowLen int
	for _, btn := range btns {
		if maxTextLen < rowLen+1+len(btn.Text) {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			rowLen = 0
			row = []telego.InlineKeyboardButton{}
		}
		rowLen += 1 + len(btn.Text)
		row = append(row, btn)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}

func (b *Bot) getButtonData(table, path string) (*db.Content, error) {
	return b.db.Select(table, path)
}
