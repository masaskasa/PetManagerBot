package main

import (
	"PetManagerBot/clients/telegram"
	eventConsumer "PetManagerBot/consumer/event-consumer"
	events "PetManagerBot/events/telegram"
	"PetManagerBot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	storagePath = "/home/user/storage_pet_manager_bot/storage.db"
	tgHost      = "api.telegram.org"
	batchSize   = 100
)

func main() {

	client := telegram.NewClient(tgHost, mustToken())

	storage, err := sqlite.NewSqliteDB(storagePath)
	if err != nil {

		log.Fatalf("can't connect to storage: %s", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: %s", err)
	}

	fetcher := events.NewFetcher(client)
	processor := events.NewProcessor(client, storage)

	consumer := eventConsumer.NewConsumer(fetcher, processor, batchSize)

	go processor.Sessions.CleanOldSessions()
	consumer.Start()
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
