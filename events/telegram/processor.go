package telegram

import (
	"PetManagerBot/clients/telegram"
	eventsPack "PetManagerBot/events"
	"PetManagerBot/handler"
	"errors"
	"log/slog"
)

type Processor struct {
	tgClient *telegram.Client
}

func NewProcessor(tgClient *telegram.Client) *Processor {
	return &Processor{
		tgClient: tgClient,
	}
}

var (
	ErrUnknownEvent = errors.New("can't process: unknown event type")
)

func (processor *Processor) Process(event eventsPack.Event) error {

	switch event.Type {
	case eventsPack.Message:
		return processor.processMessage(event)
	default:
		return ErrUnknownEvent
	}
}

func (processor *Processor) processMessage(event eventsPack.Event) error {

	if err := handler.HandleEvent(event); err != nil {
		slog.Error("processMessage: can't handle message", err)
		return err
	}

	return nil
}
