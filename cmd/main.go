package main

import (
	"log"
	"technician_bot/cmd/bot"
	"technician_bot/cmd/db"
	"technician_bot/cmd/xmlToDB"
	//"technician_bot/database"
)

func main() {
	dataBase := db.Init()
	//database.ConnectDb()
	xmlToDB.XMLToDB("kia", dataBase)

	b := bot.Init(dataBase)
	log.Println("program start")
	b.Start()
	defer b.Stop()
}
