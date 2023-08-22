package bot

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func (b *Bot) sendMenu(chatID telego.ChatID, answerForQueryID string, keyboard *telego.InlineKeyboardMarkup, text string) {
	message := telegoutil.Message(
		chatID,
		text,
	)
	if keyboard != nil {
		message.WithReplyMarkup(keyboard)
	}
	_, err := b.telegram.SendMessage(message)
	if err != nil {
		b.logger.Errorf(err.Error())
	}
	if answerForQueryID != "" {
		err = b.telegram.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{CallbackQueryID: answerForQueryID})
		if err != nil {
			b.logger.Errorf(err.Error())
		}
	}

}

// 50 символов в строке
// 2 символа между двумя кнопками
func (b *Bot) btnsOptimalPlacement(btns []telego.InlineKeyboardButton) *telego.InlineKeyboardMarkup {
	if len(btns) == 0 {
		return nil
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
