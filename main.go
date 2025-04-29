package main

import (
	"PetManagerBot/clients/telegram"
	"time"
)

const tgHost = "api.telegram.org"

func main() {
	client := telegram.New(tgHost, "top secret")

	offset := 0

	for {
		updates, _ := client.Updates(offset, 100)
		if len(updates) != 0 {
			offset = updates[len(updates)-1].ID + 1
			for _, upd := range updates {
				println(upd.Message.Text)
			}
			time.Sleep(1 * time.Second)
		}
	}
}
