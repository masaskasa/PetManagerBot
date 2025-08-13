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

	user, err := handler.session.GetObject(userName)
	if err != nil {
		return nil, err
	}

	var petsList map[uuid.UUID]*models.Pet

	list, err := handler.session.GetObject(userPets)
	if err != nil {
		petsList, err = handler.storage.GetPetsList(context.Background(), user.(string))
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

	if handler.answerCallbackQuery == nil {
		result := handler.showPetsList(msgInvalidPet)
		return result
	}

	answer, err := handler.session.GetObject(callbackQueryData)
	if err != nil {
		return err
	}

	petID, err := uuid.Parse(answer.(string))

	petsList, err := handler.session.GetObject(userPets)
	if err != nil {
		return err
	}

	pets := petsList.(map[uuid.UUID]*models.Pet)

	pet, ok := pets[petID]
	if !ok {
		if _, err := handler.answerCallbackQuery("", false); err != nil {
			slog.Error("showPetCard: answerCallbackQuery:", err)
			return err
		}
		result := handler.showPetsList(msgInvalidPet)
		return result
	}

	petCache, err := handler.session.GetObject(pet.ID.String())
	if err != nil {
		pet, err = handler.storage.GetPet(context.Background(), pet.ID)
		if err != nil {
			return err
		}
		handler.session.UpdateObject(pet.ID.String(), pet)
	} else {
		pet = petCache.(*models.Pet)
	}

	if _, err := handler.answerCallbackQuery("", false); err != nil {
		slog.Error("showPetCard: answerCallbackQuery:", err)
		return err
	}

	_, result := handler.sendMessage(pet.String())

	handler.session.setState(ready)
	handler.session.setScenario(none)
	handler.session.deleteTempObjects(messageText, callbackQueryData)

	return result
}
