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
