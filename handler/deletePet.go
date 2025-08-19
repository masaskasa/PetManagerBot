package handler

import (
	"PetManagerBot/handler/models"
	"context"
	"log/slog"
)

func (handler *Handler) startDeletePetScenario() error {

	handler.session.setScenario(deletePetCommand)
	handler.session.setState(start)

	return handler.showPetsList(msgPickPet)
}

func (handler *Handler) deletePet() error {

	switch handler.session.state {
	case start:
		return handler.askDeleteConfirmation()
	case pickPet:
		return handler.setReadyDeletePet()
	default:
		return nil
	}
}

func (handler *Handler) askDeleteConfirmation() error {

	pet, err := handler.determinePet()
	if err != nil {
		return err
	}

	handler.session.UpdateObject(deletePetCard, pet)
	handler.session.setState(pickPet)

	_, result := handler.sendMessageKeyboard(msgConfirmDeletePet+pet.String(), *handler.confirmationButtons())
	return result
}

func (handler *Handler) setReadyDeletePet() error {

	pet, err := handler.session.GetObject(deletePetCard)
	if err != nil {
		return err
	}

	deletePet, ok := pet.(*models.Pet)
	if !ok {
		slog.Error("setReadyDeletePet: type assertion problem: expected pet, get:", pet)
		return ErrExpectedAnotherType
	}

	confirmation, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return err
	}

	var result error

	switch confirmation {
	case "confirm":
		if err := handler.storage.Remove(context.Background(), deletePet.ID); err != nil {
			return err
		}
		if _, err := handler.answerCallbackQuery("", false); err != nil {
			slog.Error("setReadyDeletePet: answerCallbackQuery:", err)
			return err
		}
		_, result = handler.sendMessage(msgReadyDeletePet)
		handler.session.deleteTempObjects(userPets, deletePet.ID.String())
	case "do_not_confirm":
		if _, err := handler.answerCallbackQuery("", false); err != nil {
			slog.Error("setReadyDeletePet: answerCallbackQuery:", err)
			return err
		}
		_, result = handler.sendMessage(msgCancelDelete)
	default:
		_, result := handler.sendMessageKeyboard(msgConfirmDeletePet+deletePet.String(), *handler.confirmationButtons())
		return result
	}

	handler.session.setState(ready)
	handler.session.setScenario(none)
	handler.session.deleteTempObjects(messageText, deletePetCard, callbackQueryData)

	return result
}
