package event_consumer

import (
	eventsPack "PetManagerBot/events"
	"log/slog"
	"time"
)

type Consumer struct {
	fetcher   eventsPack.Fetcher
	processor eventsPack.Processor
	batchSize int
}

func NewConsumer(fetcher eventsPack.Fetcher, processor eventsPack.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (consumer *Consumer) Start() {
	for {
		events, err := consumer.fetcher.Fetch(consumer.batchSize)
		if err != nil {
			slog.Error("Start: Fetch:", err.Error())
			continue
		}

		if len(events) == 0 {
			time.Sleep(3 * time.Second)
			continue
		}

		consumer.handleEvents(events)
	}
}

func (consumer *Consumer) handleEvents(events []eventsPack.Event) {

	slog.Info("Start handle events")

	for _, event := range events {
		slog.Info("handleEvents: got event", event.Text)

		if err := consumer.processor.Process(&event); err != nil {
			slog.Error("handleEvents: can't handle event:", err)
			continue
		}
	}
}
