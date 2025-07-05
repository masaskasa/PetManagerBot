package handler

import (
	"PetManagerBot/clients/telegram"
	"log/slog"
	"strings"
)

const (
	chatID      = "ChatID"
	userName    = "UserName"
	messageText = "MessageText"
)

func Handle(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	if session.scenario == none {
		err := doCommand(session, sendMessage)
		if err != nil {
			slog.Error("Handle: can't do command", err)
			return err
		}
	} else {
		//TODO pass to command handler
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
	default:
		return unknownMsg(sendMessage)
	}
}

func helloMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgHello)
	return result
}

func unknownMsg(sendMessage func(string) (telegram.Message, error)) error {
	_, result := sendMessage(msgUnknownCommand)
	return result
}
