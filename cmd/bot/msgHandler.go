package bot

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"strings"
	"technician_bot/cmd/xmlToDB"
	"technician_bot/database"
)

func (b *Bot) msgHandler(msg *telego.Message) {
	var keyboard *telego.InlineKeyboardMarkup
	var text string

	// не красивое быстрое решение
	if msg.Caption != "" {
		msg.Text = msg.Caption
	}

	if msg.Text == "/start" {
		keyboard, text = b.startMsgParams()
	} else if len(msg.Text) != 0 && msg.Text[0] == '/' {
		text = b.rootCommand(msg)
	} else {
		return
	}

	b.sendMenu(telegoutil.ID(msg.Chat.ID), "", keyboard, text)
}

func (b *Bot) verify(username string) bool {
	_, ok := b.roots[username]
	return ok
}

func (b *Bot) rootCommand(msg *telego.Message) string {
	tmp := strings.Split(msg.Text, " ")

	command := tmp[0]
	data := ""
	if len(tmp) > 1 {
		data = tmp[1]
	}

	if command == "/яАдмин" {
		b.roots[msg.From.Username] = struct{}{}
		return "успешно"
	}

	if !b.verify(msg.From.Username) {
		return "Нет такой команды"
	}

	if command == "/добавитьМашину" {
		err := b.fileHandler(msg.Document)
		if err != nil {
			b.logger.Errorf(err.Error())
			return err.Error()
		}
	} else if command == "/удалитьМашину" {
		err := database.DropTable(data)
		if err != nil {
			b.logger.Errorf(err.Error())
			return err.Error()
		}
	} else {
		return "Нет такой команды"
	}
	return "успешно"
}

func (b *Bot) fileHandler(doc *telego.Document) error {
	splitFileName := strings.Split(doc.FileName, ".")

	if len(splitFileName) != 2 {
		return errors.New("incorrect file name")
	}

	tableName := splitFileName[0]
	fileType := splitFileName[1]

	if fileType != "xml" {
		return errors.New("incorrect file name")
	}

	file, err := b.telegram.GetFile(&telego.GetFileParams{FileID: doc.FileID})
	if err != nil {
		return err
	}

	fileData, err := telegoutil.DownloadFile(b.telegram.FileDownloadURL(file.FilePath))
	if err != nil {
		return err
	}

	err = xmlToDB.ByteToDB(fileData, tableName)

	return err
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

	tablesName, err := database.GetTables()
	if err != nil {
		b.logger.Error(err)
	}

	for _, name := range tablesName {
		key, err := database.GetKey(name, "start")
		if err != nil {
			b.logger.Errorf("table or first row %v not found", name)
			continue
		}
		tables[name] = fmt.Sprintf("%v@%v", name, key)
	}

	return tables
}
