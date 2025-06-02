package telegram

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/events"
)

type Fetcher struct {
	tgClient *telegram.Client
	offset   int
}

func NewFetcher(tgClient *telegram.Client) *Fetcher {
	return &Fetcher{
		tgClient: tgClient,
	}
}

func (fetcher *Fetcher) Fetch(limit int) ([]events.Event, error) {
	// TODO
	return nil, nil
}
