package telegram

import (
	"PetManagerBot/clients/telegram"
	eventsPack "PetManagerBot/events"
	"PetManagerBot/handler"
	"PetManagerBot/storage"
	"errors"
	"log/slog"
)

type ProcessorImpl struct {
	TgClient *telegram.Client
	Sessions handler.SessionsMap
	Storage  storage.Storage
}

func NewProcessor(tgClient *telegram.Client, storage storage.Storage) *ProcessorImpl {
	return &ProcessorImpl{
		TgClient: tgClient,
		Sessions: handler.NewSessionsMap(),
		Storage:  storage,
	}
}

var (
	ErrUnknownEvent    = errors.New("can't process: unknown event type")
	ErrUnknownMetaType = errors.New("can't process: unknown metaMessage type")
)

const (
	chatID            = "ChatID"
	userName          = "UserName"
	messageText       = "MessageText"
	callbackQueryID   = "CallbackQueryID"
	callbackQueryData = "CallbackQueryData"
)

func (processor *ProcessorImpl) Process(event *eventsPack.Event) error {

	switch event.Type {
	case eventsPack.Message:
		return processor.processMessage(event)
	case eventsPack.CallbackQuery:
		return processor.processCallbackQuery(event)
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

	sendMessageKeyboard, err := newMessageSenderKeyboard(session, processor.TgClient)
	if err != nil {
		slog.Error("processMessage: can't create message sender keyboard", err)
		return err
	}

	if err := handler.NewHandler(session, processor.Storage, sendMessage, sendMessageKeyboard, nil).Handle(); err != nil {
		slog.Error("processMessage: can't handle message", err)
		return err
	}

	return nil
}

func (processor *ProcessorImpl) processCallbackQuery(event *eventsPack.Event) error {

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

	sendMessageKeyboard, err := newMessageSenderKeyboard(session, processor.TgClient)
	if err != nil {
		slog.Error("processMessage: can't create message sender keyboard", err)
		return err
	}

	answerCallBackQuery, err := newAnswererCallbackQuery(session, processor.TgClient)
	if err != nil {
		slog.Error("processMessage: can't create callback query answerer", err)
		return err
	}

	if err := handler.NewHandler(session, processor.Storage, sendMessage, sendMessageKeyboard, answerCallBackQuery).Handle(); err != nil {
		slog.Error("processMessage: can't handle message", err)
		return err
	}

	return nil
}

func (processor *ProcessorImpl) prepareSession(event *eventsPack.Event) (*handler.Session, error) {

	switch event.Type {

	case eventsPack.Message:
		meta, err := getMetaMessage(event)
		if err != nil {
			return nil, err
		}

		session := processor.Sessions.GetSession(meta.UserName)
		session.UpdateObject(userName, meta.UserName)
		session.UpdateObject(chatID, meta.ChatID)
		session.UpdateObject(messageText, event.Text)

		return session, nil

	case eventsPack.CallbackQuery:
		meta, err := getMetaCallbackQuery(event)
		if err != nil {
			return nil, err
		}

		session := processor.Sessions.GetSession(meta.UserName)
		session.UpdateObject(callbackQueryID, meta.ID)
		session.UpdateObject(callbackQueryData, meta.Data)

		return session, nil
	default:
	}

	return nil, nil
}

func newMessageSender(session *handler.Session, tgClient *telegram.Client) (func(string) (telegram.Message, error), error) {

	chatID, err := session.GetInt(chatID)
	if err != nil {
		return nil, err
	}

	return func(msg string) (telegram.Message, error) {
		return tgClient.SendMessage(chatID, msg, telegram.InlineKeyboardMarkup{})
	}, nil
}

func newMessageSenderKeyboard(session *handler.Session, tgClient *telegram.Client) (func(string, telegram.InlineKeyboardMarkup) (telegram.Message, error), error) {

	chatID, err := session.GetInt(chatID)
	if err != nil {
		return nil, err
	}

	return func(msg string, keyboard telegram.InlineKeyboardMarkup) (telegram.Message, error) {
		return tgClient.SendMessage(chatID, msg, keyboard)
	}, nil
}

func newAnswererCallbackQuery(session *handler.Session, tgClient *telegram.Client) (func(string, bool) (telegram.Message, error), error) {

	callbackQueryID, err := session.GetString(callbackQueryID)
	if err != nil {
		return nil, err
	}

	return func(text string, showAlert bool) (telegram.Message, error) {
		return tgClient.AnswerCallbackQuery(callbackQueryID, text, showAlert)
	}, nil
}

func getMetaMessage(event *eventsPack.Event) (metaMessage, error) {

	result, ok := event.Meta.(metaMessage)
	if ok {
		return result, nil
	}

	return metaMessage{}, ErrUnknownMetaType
}

func getMetaCallbackQuery(event *eventsPack.Event) (metaCallbackQuery, error) {

	result, ok := event.Meta.(metaCallbackQuery)
	if ok {
		return result, nil
	}

	return metaCallbackQuery{}, ErrUnknownMetaType
}
