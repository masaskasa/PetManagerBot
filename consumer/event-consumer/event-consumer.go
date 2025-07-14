package event_consumer

import (
	eventsPack "PetManagerBot/events"
	"log/slog"
	"time"
)

type ConsumerImpl struct {
	fetcher   eventsPack.Fetcher
	processor eventsPack.Processor
	batchSize int
}

func NewConsumer(fetcher eventsPack.Fetcher, processor eventsPack.Processor, batchSize int) *ConsumerImpl {
	return &ConsumerImpl{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

const noEventsDelay = 1 * time.Second

func (consumer *ConsumerImpl) Start() {

	for {
		events, err := consumer.fetcher.Fetch(consumer.batchSize)
		if err != nil {
			slog.Error("Start: Fetch:", err.Error())
			continue
		}

		slog.Info("Fetch events:", events)

		if len(events) == 0 {
			time.Sleep(noEventsDelay)
			continue
		}

		consumer.handleEvents(events)
	}
}

func (consumer *ConsumerImpl) handleEvents(events []eventsPack.Event) {

	slog.Info("Start handle events")

	for _, event := range events {
		slog.Info("handleEvents: got event", event.Text)

		if err := consumer.processor.Process(&event); err != nil {
			slog.Error("handleEvents: can't handle event:", err)
			continue
		}
	}
}
