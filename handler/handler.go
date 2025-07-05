package handler

import (
	"PetManagerBot/clients/telegram"
	storagePack "PetManagerBot/storage"
	"log/slog"
	"strings"
)

const (
	messageText = "MessageText"
	userName    = "UserName"
	pet         = "Pet"
	species     = "Species"
	breed       = "Breeds"
)

func Handle(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	switch session.scenario {

	case none:
		err := doCommand(session, sendMessage)
		if err != nil {
			slog.Error("Handle: can't do command", err)
			return err
		}

	case createPetCommand:
		return createPet(session, sendMessage, storage)

	}

	return nil
}

func doCommand(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	text, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	command := strings.TrimSpace(text.(string))

	switch command {
	case startCommand:
		return helloMsg(sendMessage)
	case helpCommand:
		return helpMsg(sendMessage)
	case createPetCommand:
		doCreatePetScenario(session)
		return nameMsg(sendMessage)
	default:
		return unknownMsg(sendMessage)
	}
}

func helloMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgHello)
	return result
}

func nameMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgAskName)
	return result
}

func helpMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgHelp)
	return result
}

func unknownMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgUnknownCommand)
	return result
}
