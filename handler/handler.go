package handler

import (
	"PetManagerBot/clients/telegram"
	storagePack "PetManagerBot/storage"
	"fmt"
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

	isBreak, err := breakScenario(session, sendMessage)
	if err != nil {
		return err
	}

	if isBreak {
		return nil
	}

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
		startCreatePetScenario(session)
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
	_, result := sendMessage(msgHowToBegin)
	return result
}

func unknownMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgUnknownCommand)
	return result
}

func breakScenario(session *Session, sendMessage func(string) (telegram.Message, error)) (bool, error) {

	text, err := session.GetObject(messageText)
	if err != nil {
		return false, err
	}

	command := strings.TrimSpace(text.(string))

	if command == breakCommand {
		return true, breakMsg(session, sendMessage)
	}

	return false, nil
}

func breakMsg(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	breakScenario := session.scenario

	session.setState(ready)
	session.setScenario(none)

	_, result := sendMessage(msgBreakCommand + fmt.Sprint(msgTryAgain, breakScenario))
	return result
}
