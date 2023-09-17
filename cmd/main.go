package main

import (
	"technician_bot/cmd/bot"
	"technician_bot/cmd/db"
	"technician_bot/cmd/xmlToDB"
	//"technician_bot/database"
)

func main() {
	dataBase := db.Init()
	//database.ConnectDb()
	xmlToDB.XMLToDB("kia", dataBase)
	xmlToDB.XMLToDB("polo", dataBase)

	b := bot.Init(dataBase)
	b.Start()
	defer b.Stop()
}
