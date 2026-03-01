package handler

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/handler/models"
	"errors"
	"fmt"
	"log/slog"
)

func (handler *Handler) startEditPetScenario() error {

	_, result := handler.sendMessage("Sorry! Not implemented yet")
	return result

	//handler.session.setScenario(editPetCommand)
	//handler.session.setState(start)
	//
	//return handler.showPetsList(msgPickPet)
}

func (handler *Handler) editPet() error {

	switch handler.session.state {
	case start:
		return handler.showParametersList()
	case pickPet:
		return nil
	case pickParameter:
		return nil
	case parameterComplete:
		return nil
	default:
		return nil
	}
}

func (handler *Handler) showParametersList() error {

	pet, err := handler.determinePet()
	if err != nil {
		return err
	}

	pet, err = handler.getFullPetCard(pet)
	if err != nil {
		return err
	}

	handler.session.UpdateObject(editPetCard, pet)
	handler.session.setState(pickPet)

	msgPersonal := fmt.Sprintf(msgPickParameter, pet.Name)

	_, result := handler.sendMessageKeyboard(msgPersonal+pet.String(), *handler.parametersButtons())
	return result
}

func (handler *Handler) askNewValueForParameter() error {

	pet, err := handler.determineEditPet()
	if err != nil {
		return err
	}

	parameter, err := handler.determineParameter(pet)
	if err != nil {
		return err
	}

	switch parameter {

	}

	handler.session.setState(pickParameter)

	return err
}

func (handler *Handler) parametersButtons() *telegram.InlineKeyboardMarkup {

	buttons := telegram.NewInlineKeyboardMarkup()

	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "✨ Name", CallbackData: models.EditName})
	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "🐾 Species", CallbackData: models.EditSpecies})
	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "🏷️ Breed", CallbackData: models.EditBreed})
	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "⚧️ Gender", CallbackData: models.EditGender})
	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "🆔 Animal ID", CallbackData: models.EditAnimalID})
	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "🔍 Special Signs", CallbackData: models.EditSpecialSigns})

	return &buttons
}

func (handler *Handler) determineParameter(pet *models.Pet) (string, error) {

	if handler.answerCallbackQuery == nil {
		_, result := handler.sendMessageKeyboard(msgInvalidParameter+pet.String(), *handler.parametersButtons())
		return "", errors.Join(errors.New("invalid parameter"), result)
	}

	answer, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func (handler *Handler) determineEditPet() (*models.Pet, error) {

	pet, err := handler.session.GetObject(editPetCard)
	if err != nil {
		return nil, err
	}

	newPet, ok := pet.(*models.Pet)
	if !ok {
		slog.Error("determineEditPet: type assertion problem: expected pet, get:", pet)
		return nil, ErrExpectedAnotherType
	}

	return newPet, nil
}
