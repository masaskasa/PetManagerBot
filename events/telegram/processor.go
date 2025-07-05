package telegram

import (
	"PetManagerBot/clients/telegram"
	eventsPack "PetManagerBot/events"
	"PetManagerBot/handler"
	"errors"
	"log/slog"
)

type ProcessorImpl struct {
	TgClient *telegram.Client
	Sessions handler.SessionsMap
}

func NewProcessor(tgClient *telegram.Client) *ProcessorImpl {
	return &ProcessorImpl{
		TgClient: tgClient,
		Sessions: handler.NewSessionsMap(),
	}
}

var (
	ErrUnknownEvent    = errors.New("can't process: unknown event type")
	ErrUnknownMetaType = errors.New("can't process: unknown meta type")
)

const (
	chatID      = "ChatID"
	userName    = "UserName"
	messageText = "MessageText"
)

func (processor *ProcessorImpl) Process(event *eventsPack.Event) error {

	switch event.Type {
	case eventsPack.Message:
		return processor.processMessage(event)
	default:
		return ErrUnknownEvent
	}
}

func (processor *ProcessorImpl) processMessage(event *eventsPack.Event) error {

	session, err := processor.prepareSession(event)
	if err != nil {
		slog.Error("processMessage: can't prepare session", err)
		return err
	}

	sendMessage, err := newMessageSender(session, processor.TgClient)
	if err != nil {
		slog.Error("processMessage: can't create message sender", err)
		return err
	}

	if err := handler.Handle(session, sendMessage); err != nil {
		slog.Error("processMessage: can't handle message", err)
		return err
	}

	return nil
}

func (processor *ProcessorImpl) prepareSession(event *eventsPack.Event) (*handler.Session, error) {

	meta, err := getMeta(event)
	if err != nil {
		return nil, err
	}

	session := processor.Sessions.GetSession(meta.UserName, meta.ChatID)
	session.UpdateObject(userName, meta.UserName)
	session.UpdateObject(chatID, meta.ChatID)
	session.UpdateObject(messageText, event.Text)

	return session, nil
}

func newMessageSender(session *handler.Session, tgClient *telegram.Client) (func(string) (telegram.Message, error), error) {

	chat, err := session.GetObject(chatID)
	if err != nil {
		return nil, err
	}

	chatID := chat.(int)

	return func(msg string) (telegram.Message, error) {
		return tgClient.SendMessage(chatID, msg)
	}, nil
}

func getMeta(event *eventsPack.Event) (Meta, error) {

	result, ok := event.Meta.(Meta)
	if ok {
		return result, nil
	} else {
		return Meta{}, ErrUnknownMetaType
	}
}
