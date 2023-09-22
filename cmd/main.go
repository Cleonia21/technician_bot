package main

import (
	"technician_bot/cmd/bot"
	"technician_bot/database"
)

func main() {
	database.ConnectDB()
	defer database.Close()

	//xmlToDB.XMLToDB("kia")
	//xmlToDB.XMLToDB("polo")

	b := bot.Init()
	b.Start()
	defer b.Stop()
}
