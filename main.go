package main

import "main/bot"

func main() {
	b := bot.Init()
	b.Start()
	defer b.Stop()
}
