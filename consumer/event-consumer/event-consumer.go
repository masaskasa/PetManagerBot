package event_consumer

import (
	eventsPack "PetManagerBot/events"
	"log/slog"
	"time"
)

type Consumer struct {
	fetcher   eventsPack.IFetcher
	processor eventsPack.IProcessor
	batchSize int
}

func New(fetcher eventsPack.IFetcher, processor eventsPack.IProcessor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (consumer *Consumer) Start() error {
	for {
		events, err := consumer.fetcher.Fetch(consumer.batchSize)
		if err != nil {
			slog.Error("Start: Fetch:", err.Error())
			continue
		}

		if len(events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		//if err := consumer.handleEvents(events); err != nil {
		//	slog.Error("Start: handleEvents:", err)
		//	continue
		//}
	}
}

func (consumer *Consumer) handleEvents(events []eventsPack.Event) {
	for _, event := range events {
		slog.Info("handleEvents: got event", event.Text)

		if err := consumer.processor.Process(event); err != nil {
			slog.Error("handleEvents: can't handle event:", err)
			continue
		}
	}
}
