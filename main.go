package main

import (
	"main/bot"
	"main/db"
)

func main() {
	dataBase := db.Init()

	//xmlToDB.XMLToDB("kia", dataBase)

	b := bot.Init(dataBase)
	b.Start()
	defer b.Stop()
}

// кнопки назад и главное меню

//ГОТОВО
// всегда одинаковый порядок кнопок в меню +
// нормальное заполнение пространства кнопками в меню +
