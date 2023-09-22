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

	if msg.Text == "/start" {
		keyboard, text = b.startMsgParams()
	} else if msg.Document != nil {
		err := b.fileHandler(msg.Document)
		if err != nil {
			b.logger.Errorf(err.Error())
		}
		keyboard, text = b.startMsgParams()
	} else {
		return
	}

	b.sendMenu(telegoutil.ID(msg.Chat.ID), "", keyboard, text)
}

func (b *Bot) fileHandler(doc *telego.Document) error {
	splitFileName := strings.Split(doc.FileName, "@")
	if len(splitFileName) != 2 || splitFileName[0] != "pass" {
		return errors.New("incorrect file name")
	}

	splitTableName := strings.Split(splitFileName[1], ".")
	if len(splitTableName) != 2 || splitTableName[1] != "xml" {
		return errors.New("incorrect file name")
	}
	tableName := splitTableName[0]

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
