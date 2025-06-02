package telegram

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/events"
)

type Processor struct {
	tgClient *telegram.Client
}

func NewProcessor(tgClient *telegram.Client) *Processor {
	return &Processor{
		tgClient: tgClient,
	}
}

func (processor *Processor) Process(event events.Event) error {
	switch event.Type {
	// TODO
	}
	return nil
}
