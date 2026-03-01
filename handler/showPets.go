package handler

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/handler/models"
	storagePack "PetManagerBot/storage"
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"sort"
)

func (handler *Handler) startShowPetScenario() error {

	handler.session.setScenario(showPetCommand)
	handler.session.setState(start)

	return handler.showPetsList(msgPickPet)
}

func (handler *Handler) showPetsList(message string) error {

	buttons, err := handler.petsButtons()
	if err != nil {
		return err
	}

	_, result := handler.sendMessageKeyboard(message, *buttons)
	return result
}

func (handler *Handler) petsButtons() (*telegram.InlineKeyboardMarkup, error) {

	user, err := handler.session.GetString(userName)
	if err != nil {
		return nil, err
	}

	var petsList map[uuid.UUID]*models.Pet

	list, err := handler.session.GetObject(userPets)
	if err != nil {

		if errors.Is(err, ErrObjectNotExists) {
			petsList, err = handler.storage.GetPetsList(context.Background(), user)
			if err != nil {
				if errors.Is(err, storagePack.ErrNoSavedPets) {
					_, result := handler.sendMessage(msgNoSavedPets)
					handler.session.setState(ready)
					handler.session.setScenario(none)
					handler.session.deleteTempObjects(messageText, callbackQueryData)
					return nil, errors.Join(storagePack.ErrNoSavedPets, result)
				}
				return nil, err
			}
			handler.session.UpdateObject(userPets, petsList)
		} else {
			return nil, err
		}

	} else {

		_, ok := list.(map[uuid.UUID]*models.Pet)
		if !ok {
			slog.Error("petsButtons: type assertion problem: expected pets, get:", list)
			return nil, ErrExpectedAnotherType
		}

		petsList = list.(map[uuid.UUID]*models.Pet)
	}

	sortPetsList := make([]*models.Pet, 0, len(petsList))
	for _, pet := range petsList {
		sortPetsList = append(sortPetsList, pet)
	}
	sort.Slice(sortPetsList, func(i, j int) bool {
		return sortPetsList[i].Name < sortPetsList[j].Name
	})

	buttons := telegram.NewInlineKeyboardMarkup()

	for _, pet := range sortPetsList {
		buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: pet.Name, CallbackData: pet.ID.String()})
	}

	return &buttons, nil
}

func (handler *Handler) showPetCard() error {

	pet, err := handler.determinePet()
	if err != nil {
		return err
	}

	pet, err = handler.getFullPetCard(pet)
	if err != nil {
		return err
	}

	_, result := handler.sendMessage(pet.String())

	handler.session.setState(ready)
	handler.session.setScenario(none)
	handler.session.deleteTempObjects(messageText, callbackQueryData)

	return result
}

func (handler *Handler) determinePet() (*models.Pet, error) {

	if handler.answerCallbackQuery == nil {
		result := handler.showPetsList(msgInvalidPet)
		return nil, errors.Join(errors.New("invalid pet"), result)
	}

	answer, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return nil, err
	}

	petID, err := uuid.Parse(answer)

	petsList, err := handler.session.GetObject(userPets)
	if err != nil {
		return nil, err
	}

	pets, ok := petsList.(map[uuid.UUID]*models.Pet)
	if !ok {
		slog.Error("determinePet: type assertion problem: expected pets, get:", petsList)
		return nil, ErrExpectedAnotherType
	}

	pet, ok := pets[petID]
	if !ok {
		if _, err := handler.answerCallbackQuery("", false); err != nil {
			slog.Error("determinePet: answerCallbackQuery:", err)
			return nil, err
		}
		result := handler.showPetsList(msgInvalidPet)
		return nil, errors.Join(errors.New("invalid pet"), result)
	}

	if _, err := handler.answerCallbackQuery("", false); err != nil {
		slog.Error("determinePet: answerCallbackQuery:", err)
		return nil, err
	}

	return pet, nil
}

func (handler *Handler) getFullPetCard(pet *models.Pet) (*models.Pet, error) {

	petCache, err := handler.session.GetObject(pet.ID.String())
	if err != nil {

		if errors.Is(err, ErrObjectNotExists) {
			pet, err = handler.storage.GetPet(context.Background(), pet.ID)
			if err != nil {
				return nil, err
			}
			handler.session.UpdateObject(pet.ID.String(), pet)
		} else {
			return nil, err
		}

	} else {

		_, ok := petCache.(*models.Pet)
		if !ok {
			slog.Error("showPetCard: type assertion problem: expected pet, get:", petCache)
			return nil, ErrExpectedAnotherType
		}

		pet = petCache.(*models.Pet)
	}
	return pet, nil
}
