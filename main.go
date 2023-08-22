package main

import (
	"main/bot"
	"main/db"
	"main/xmlToDB"
)

func main() {
	dataBase := db.Init()

	xmlToDB.XMLToDB("kia", dataBase)

	b := bot.Init(dataBase)
	b.Start()
	defer b.Stop()
}
