package main

import (
	"PetManagerBot/clients/telegram"
	"time"
)

const tgHost = "https://api.telegram.org"

func main() {
	client := telegram.New(tgHost, "")

	offset := 0

	for {
		updates, _ := client.Updates(offset, 100)
		if len(updates) != 0 {
			offset = updates[len(updates)-1].ID + 1
			println(updates[0].Message.Text)
		}
		time.Sleep(1 * time.Second)
	}
}
