package handler

import (
	"PetManagerBot/clients/telegram"
	storagePack "PetManagerBot/storage"
	"fmt"
	"log/slog"
	"strings"
)

const (
	messageText       = "MessageText"
	userName          = "UserName"
	newPetCard        = "NewPet"
	species           = "Species"
	breeds            = "Breeds"
	callbackQueryData = "CallbackQueryData"
	userPets          = "UserPets"
)

type Handler struct {
	session             *Session
	storage             storagePack.Storage
	sendMessage         func(string) (telegram.Message, error)
	sendMessageKeyboard func(string, telegram.InlineKeyboardMarkup) (telegram.Message, error)
	answerCallbackQuery func(string, bool) (telegram.Message, error)
}

func NewHandler(session *Session, storage storagePack.Storage, sendMessage func(string) (telegram.Message, error), sendMessageKeyboard func(string, telegram.InlineKeyboardMarkup) (telegram.Message, error), answerCallbackQuery func(string, bool) (telegram.Message, error)) *Handler {
	return &Handler{
		session:             session,
		storage:             storage,
		sendMessage:         sendMessage,
		sendMessageKeyboard: sendMessageKeyboard,
		answerCallbackQuery: answerCallbackQuery,
	}
}

func (handler *Handler) Handle() error {

	isBreak, err := handler.breakScenario()
	if err != nil {
		return err
	}

	if isBreak {
		return nil
	}

	switch handler.session.scenario {

	case none:
		err := handler.doCommand()
		if err != nil {
			slog.Error("Handle: can't do command", err)
			return err
		}

	case createPetCommand:
		return handler.createPet()

	case showPetCommand:
		return handler.showPetCard()
	}

	return nil
}

func (handler *Handler) doCommand() error {

	text, err := handler.session.GetString(messageText)
	if err != nil {
		return err
	}

	command := strings.TrimSpace(text)

	switch command {
	case startCommand:
		return handler.helloMsg()
	case helpCommand:
		return handler.helpMsg()
	case createPetCommand:
		handler.startCreatePetScenario()
		return handler.nameMsg()
	case showPetCommand:
		return handler.startShowPetScenario()
	default:
		return handler.unknownMsg()
	}
}

func (handler *Handler) helloMsg() error {
	_, result := handler.sendMessage(msgHello)
	return result
}

func (handler *Handler) nameMsg() error {
	_, result := handler.sendMessage(msgAskName)
	return result
}

func (handler *Handler) helpMsg() error {
	_, result := handler.sendMessage(msgHowToBegin)
	return result
}

func (handler *Handler) unknownMsg() error {
	_, result := handler.sendMessage(msgUnknownCommand)
	return result
}

func (handler *Handler) breakScenario() (bool, error) {

	text, err := handler.session.GetString(messageText)
	if err != nil {
		return false, err
	}

	command := strings.TrimSpace(text)

	if command == breakCommand {
		return true, handler.breakMsg()
	}

	return false, nil
}

func (handler *Handler) breakMsg() error {

	breakScenario := handler.session.scenario

	if breakScenario == none {
		_, result := handler.sendMessage(msgNeedlessBreakCommand)
		return result
	}

	handler.session.setState(ready)
	handler.session.setScenario(none)

	_, result := handler.sendMessage(msgBreakCommand + fmt.Sprint(msgTryAgain, breakScenario))
	return result
}
