package main

import (
	"technician_bot/cmd/bot"
	"technician_bot/cmd/xmlToDB"
	"technician_bot/database"
)

func main() {
	database.ConnectDB()
	//dataBase := database.Init()
	//database.ConnectDb()
	xmlToDB.XMLToDB("kia")
	xmlToDB.XMLToDB("polo")

	b := bot.Init()
	b.Start()
	defer b.Stop()
	database.Close()
}
