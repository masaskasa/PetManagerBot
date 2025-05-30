package telegram

import (
	"PetManagerBot/clients/telegram"
	"log/slog"
)

type Fetcher struct {
	tgClient *telegram.Client
	offset   int
}

func New(tgClient *telegram.Client) *Fetcher {
	return &Fetcher{
		tgClient: tgClient,
	}
}

func (fetcher *Fetcher) Fetch(limit int) ([]telegram.Update, error) {

	updates, err := fetcher.tgClient.GetUpdates(fetcher.offset, limit)
	if err != nil {
		slog.Error("Fetch: can't get updates", err.Error())
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	fetcher.offset = updates[len(updates)-1].ID + 1

	return updates, nil
}
