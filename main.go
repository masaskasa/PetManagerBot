package main

import (
	"PetManagerBot/clients/telegram"
	"flag"
	"log"
	"log/slog"
	"time"
)

const tgHost = "api.telegram.org"

func main() {
	client := telegram.NewClient(tgHost, mustToken())

	offset := 0
	for {
		updates, _ := client.Updates(offset, 100)
		if len(updates) != 0 {
			offset = updates[len(updates)-1].ID + 1
			for _, upd := range updates {
				slog.Debug(upd.Message.Text)
			}
			time.Sleep(1 * time.Second)
		} else {
			continue
		}

		// test POST request
		//data, err := json.Marshal(telegram.NewTextMessage(updates[0].Message.Chat.ID, "hello handsome"))
		//if err != nil {
		//	log.Fatal("error marshal text message")
		//}
		//
		//body, err := client.PostRequest(telegram.SendMessage, data)
		//if err != nil {
		//	log.Fatal("error marshal text message")
		//}
		//
		//slog.Info("response body:", body)
	}
}

/*
Return a token for client.
Token parse from command line argument with flag -token-for-bot.
If argument empty print error and call to os.Exit(1)
*/
func mustToken() string {
	token := flag.String("token-for-bot", "", "set token for telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
