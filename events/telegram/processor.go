package telegram

import (
	"PetManagerBot/clients/telegram"
	eventsPack "PetManagerBot/events"
	"PetManagerBot/handler"
	"errors"
	"log/slog"
)

type ProcessorImpl struct {
	tgClient *telegram.Client
	sessions handler.SessionsMap
}

func NewProcessor(tgClient *telegram.Client) *ProcessorImpl {
	return &ProcessorImpl{
		tgClient: tgClient,
		sessions: handler.NewSessionsMap(),
	}
}

var (
	ErrUnknownEvent = errors.New("can't process: unknown event type")
)

func (processor *ProcessorImpl) Process(event eventsPack.Event) error {

	switch event.Type {
	case eventsPack.Message:
		return processor.processMessage(event)
	default:
		return ErrUnknownEvent
	}
}

func (processor *ProcessorImpl) processMessage(event eventsPack.Event) error {

	if err := handler.Handle(event); err != nil {
		slog.Error("processMessage: can't handle message", err)
		return err
	}

	return nil
}
