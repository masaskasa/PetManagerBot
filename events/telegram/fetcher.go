package telegram

import (
	"PetManagerBot/clients/telegram"
	eventsPack "PetManagerBot/events"
	"log/slog"
)

type FetcherImpl struct {
	tgClient *telegram.Client
	offset   int
}

func NewFetcher(tgClient *telegram.Client) *FetcherImpl {
	return &FetcherImpl{
		tgClient: tgClient,
	}
}

func (fetcher *FetcherImpl) Fetch(limit int) ([]eventsPack.Event, error) {

	updates, err := fetcher.tgClient.GetUpdates(fetcher.offset, limit)
	if err != nil {
		slog.Error("Fetch: can't get events", err)
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	events := make([]eventsPack.Event, len(updates))

	for _, update := range updates {
		events = append(events, makeEvent(update))
	}

	slog.Info("Fetch events:", events)

	fetcher.offset = updates[len(updates)-1].ID + 1

	return nil, nil
}

func makeEvent(update telegram.Update) eventsPack.Event {

	event := eventsPack.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}

	if event.Type == eventsPack.Message {
		event.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	}

	return event
}

type Meta struct {
	ChatID   int
	UserName string
}

func fetchType(update telegram.Update) eventsPack.Type {

	if update.Message == nil {
		return eventsPack.Unknown
	}

	return eventsPack.Message
}

func fetchText(update telegram.Update) string {

	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
