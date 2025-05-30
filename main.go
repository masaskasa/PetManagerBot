package main

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const tgHost = "api.telegram.org"

func main() {
	client := telegram.NewClient(tgHost, mustToken())

	storage, err := sqlite.NewSqliteDB("/home/user/storage_pet_manager_bot/storage.db")
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: %s", err)
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
