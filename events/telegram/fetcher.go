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
		return []eventsPack.Event{}, nil
	} else {
		for _, update := range updates {
			slog.Info("Fetch: GetUpdates: CallbackQuery:", update.CallbackQuery, "Message:", update.Message)
		}
	}

	events := make([]eventsPack.Event, len(updates))

	for _, update := range updates {
		events = append(events, makeEvent(update))
	}

	fetcher.offset = updates[len(updates)-1].ID + 1

	return events, nil
}

func makeEvent(update telegram.Update) eventsPack.Event {

	event := eventsPack.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}

	switch event.Type {
	case eventsPack.Message:
		event.Meta = metaMessage{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	case eventsPack.CallbackQuery:
		event.Meta = metaCallbackQuery{
			ID:       update.CallbackQuery.ID,
			UserName: update.CallbackQuery.From.UserName,
			Data:     update.CallbackQuery.Data,
		}
	default:
	}

	return event
}

type metaMessage struct {
	ChatID   int
	UserName string
}

type metaCallbackQuery struct {
	ID       string
	UserName string
	Data     string
}

func fetchType(update telegram.Update) eventsPack.Type {

	if update.CallbackQuery != nil {
		return eventsPack.CallbackQuery
	}

	if update.Message != nil {
		return eventsPack.Message
	}

	return eventsPack.Unknown
}

func fetchText(update telegram.Update) string {

	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
