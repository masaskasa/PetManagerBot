package handler

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/handler/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
)

func (handler *Handler) startCreatePetScenario() {
	handler.session.setScenario(createPetCommand)
	handler.session.setState(start)
}

var (
	ErrUnknownSpecies      = errors.New("unknown species")
	ErrUnknownBreed        = errors.New("unknown breed")
	ErrUnknownSex          = errors.New("unknown sex")
	ErrExpectedAnotherType = errors.New("expected another type")
)

func (handler *Handler) createPet() error {

	switch handler.session.state {
	case start:
		return handler.setNameComplete()
	case nameComplete:
		return handler.setSpeciesComplete()
	case speciesComplete:
		return handler.setBreedComplete()
	case breedComplete:
		return handler.setSexComplete()
	case sexComplete:
		return handler.setAnimalIDComplete()
	case animalIDComplete:
		return handler.setSpecialSignsComplete()
	case specialSignsComplete:
		return handler.setReadyCreatePet()
	default:
		return nil
	}
}

func (handler *Handler) setNameComplete() error {

	answer, err := handler.session.GetString(messageText)
	if err != nil {
		return err
	}

	owner, err := handler.session.GetString(userName)
	if err != nil {
		return err
	}

	newPet := models.NewPet(owner)

	if err := newPet.SetName(answer); err != nil {
		return err
	}

	handler.session.UpdateObject(newPetCard, newPet)
	handler.session.setState(nameComplete)

	speciesButtons, err := handler.speciesButtons()
	if err != nil {
		return err
	}

	_, result := handler.sendMessageKeyboard(msgAskSpecies, *speciesButtons)
	return result
}

func (handler *Handler) setSpeciesComplete() error {

	if handler.answerCallbackQuery == nil {
		speciesButtons, err := handler.speciesButtons()
		if err != nil {
			return err
		}
		_, result := handler.sendMessageKeyboard(msgInvalidSpecies, *speciesButtons)
		return result
	}

	species, err := handler.determineSpecies()
	if errors.Is(err, ErrUnknownSpecies) {
		speciesButtons, err := handler.speciesButtons()
		if err != nil {
			return err
		}
		_, result := handler.sendMessageKeyboard(msgInvalidSpecies, *speciesButtons)
		return result
	} else if err != nil {
		return err
	}

	if _, err := handler.answerCallbackQuery(ntfSetSpecies+species.Name, false); err != nil {
		slog.Error("setSpeciesComplete: answerCallbackQuery:", err)
	}

	newPet, err := handler.determineNewPet()
	if err != nil {
		return err
	}

	newPet.SetSpecies(species)

	handler.session.UpdateObject(newPetCard, newPet)
	handler.session.setState(speciesComplete)

	breedButtons, err := handler.breedButtons(species.ID)
	if err != nil {
		return err
	}

	_, result := handler.sendMessageKeyboard(msgAskBreed, *breedButtons)
	return result
}

func (handler *Handler) setBreedComplete() error {

	newPet, err := handler.determineNewPet()
	if err != nil {
		return err
	}

	if handler.answerCallbackQuery == nil {
		breedButtons, err := handler.breedButtons(newPet.Species.ID)
		if err != nil {
			return err
		}
		_, result := handler.sendMessageKeyboard(msgInvalidBreed, *breedButtons)
		return result
	}

	breed, err := handler.determineBreed()
	if errors.Is(err, ErrUnknownBreed) {
		breedButtons, err := handler.breedButtons(newPet.Species.ID)
		if err != nil {
			return err
		}
		_, result := handler.sendMessageKeyboard(msgInvalidBreed, *breedButtons)
		return result
	} else if err != nil {
		return err
	}

	newPet.SetBreed(breed)

	if _, err := handler.answerCallbackQuery(ntfSetBreed+breed.Name, false); err != nil {
		slog.Error("setBreedComplete: answerCallbackQuery:", err)
	}

	handler.session.UpdateObject(newPetCard, newPet)
	handler.session.setState(breedComplete)

	_, result := handler.sendMessageKeyboard(msgAskSex, *handler.sexButtons())
	return result
}

func (handler *Handler) setSexComplete() error {

	if handler.answerCallbackQuery == nil {
		_, result := handler.sendMessageKeyboard(msgInvalidSex, *handler.sexButtons())
		return result
	}

	sex, err := handler.determineSex()
	if errors.Is(err, ErrUnknownSex) {
		_, result := handler.sendMessageKeyboard(msgInvalidSex, *handler.sexButtons())
		return result
	} else if err != nil {
		return err
	}

	newPet, err := handler.determineNewPet()
	if err != nil {
		return err
	}

	newPet.SetSex(sex)

	if _, err := handler.answerCallbackQuery(ntfSetSex+sex.String(), false); err != nil {
		slog.Error("setSexComplete: answerCallbackQuery:", err)
	}

	handler.session.UpdateObject(newPetCard, newPet)
	handler.session.setState(sexComplete)

	_, result := handler.sendMessageKeyboard(msgAskAnimalID, *handler.skipButton())
	return result
}

func (handler *Handler) setAnimalIDComplete() error {

	if handler.answerCallbackQuery != nil {
		if _, err := handler.answerCallbackQuery(ntfSkip, false); err != nil {
			slog.Error("setAnimalIDComplete: answerCallbackQuery:", err)
		}
	} else {
		answer, err := handler.session.GetString(messageText)
		if err != nil {
			return err
		}

		newPet, err := handler.determineNewPet()
		if err != nil {
			return err
		}

		animalID := strings.Trim(answer, "/")
		animalID = strings.ToLower(animalID)

		if err := newPet.SetAnimalID(answer); err != nil {
			return err
		}
		handler.session.UpdateObject(newPetCard, newPet)
	}

	handler.session.setState(animalIDComplete)

	_, result := handler.sendMessageKeyboard(msgAskSpecialSigns, *handler.skipButton())
	return result
}

func (handler *Handler) setSpecialSignsComplete() error {

	newPet, err := handler.determineNewPet()
	if err != nil {
		return err
	}

	if handler.answerCallbackQuery != nil {
		if _, err := handler.answerCallbackQuery(ntfSkip, false); err != nil {
			slog.Error("setSpecialSignsComplete: answerCallbackQuery:", err)
		}
	} else {
		answer, err := handler.session.GetString(messageText)
		if err != nil {
			return err
		}

		newPet.SetSpecialSigns(answer)

		handler.session.UpdateObject(newPetCard, newPet)
	}

	handler.session.setState(specialSignsComplete)

	_, result := handler.sendMessageKeyboard(msgConfirmCreatePet+newPet.String(), *handler.confirmationButtons())
	return result
}

func (handler *Handler) setReadyCreatePet() error {

	newPet, err := handler.determineNewPet()
	if err != nil {
		return err
	}

	confirmation, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return err
	}

	var result error

	switch confirmation {
	case "confirm":
		if err := handler.storage.Save(context.Background(), newPet); err != nil {
			return err
		}
		if _, err := handler.answerCallbackQuery("", false); err != nil {
			slog.Error("setReadyCreatePet: answerCallbackQuery:", err)
			return err
		}
		_, result = handler.sendMessage(msgReadyCreatePet)
	case "do_not_confirm":
		if _, err := handler.answerCallbackQuery("", false); err != nil {
			slog.Error("setReadyCreatePet: answerCallbackQuery:", err)
			return err
		}
		_, result = handler.sendMessage(fmt.Sprint(msgTryAgain, handler.session.scenario))
	default:
		_, result := handler.sendMessageKeyboard(msgConfirmCreatePet+newPet.String(), *handler.confirmationButtons())
		return result
	}

	handler.session.setState(ready)
	handler.session.setScenario(none)
	handler.session.deleteTempObjects(messageText, newPetCard, callbackQueryData)

	return result
}

func (handler *Handler) determineNewPet() (*models.Pet, error) {

	pet, err := handler.session.GetObject(newPetCard)
	if err != nil {
		return nil, err
	}

	newPet, ok := pet.(*models.Pet)
	if !ok {
		slog.Error("determinePet: type assertion problem: expected pet, get:", pet)
		return nil, ErrExpectedAnotherType
	}

	return newPet, nil
}

func (handler *Handler) determineSpecies() (*models.Species, error) {

	serializedID, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, ErrUnknownSpecies
	}

	list, err := handler.session.GetObject(species)
	if err != nil {
		return nil, err
	}

	speciesList, ok := list.(map[int]*models.Species)
	if !ok {
		slog.Error("determineSpecies: type assertion problem: expected species, get:", list)
		return nil, ErrExpectedAnotherType
	}

	species, exists := speciesList[id]
	if !exists {
		return nil, ErrUnknownSpecies
	}

	return species, nil
}

func (handler *Handler) determineBreed() (*models.Breed, error) {

	serializedID, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, ErrUnknownBreed
	}

	list, err := handler.session.GetObject(breeds)
	if err != nil {
		return nil, err
	}

	breedList, ok := list.(map[int]*models.Breed)
	if !ok {
		slog.Error("determineBreed: type assertion problem: expected breeds, get:", list)
		return nil, ErrExpectedAnotherType
	}

	breed, exists := breedList[id]
	if !exists {
		return nil, ErrUnknownBreed
	}

	return breed, nil
}

func (handler *Handler) determineSex() (models.Sex, error) {

	sex, err := handler.session.GetString(callbackQueryData)
	if err != nil {
		return 0, err
	}

	switch sex {
	case "female":
		return models.Female, nil
	case "male":
		return models.Male, nil
	default:
		return 0, ErrUnknownSex
	}
}

func (handler *Handler) speciesButtons() (*telegram.InlineKeyboardMarkup, error) {

	var speciesList map[int]*models.Species

	list, err := handler.session.GetObject(species)
	if err != nil {

		if errors.Is(err, ErrObjectNotExists) {
			speciesList, err = handler.storage.GetSpeciesList(context.Background())
			if err != nil {
				return nil, err
			}
			handler.session.UpdateObject(species, speciesList)
		} else {
			return nil, err
		}

	} else {

		_, ok := list.(map[int]*models.Species)
		if !ok {
			slog.Error("speciesButtons: type assertion problem: expected species, get:", list)
			return nil, ErrExpectedAnotherType
		}

		speciesList = list.(map[int]*models.Species)
	}

	sortSpeciesList := make([]*models.Species, 0, len(speciesList))
	for _, species := range speciesList {
		sortSpeciesList = append(sortSpeciesList, species)
	}
	sort.Slice(sortSpeciesList, func(i, j int) bool {
		return sortSpeciesList[i].Name < sortSpeciesList[j].Name
	})

	buttons := telegram.NewInlineKeyboardMarkup()

	for _, species := range sortSpeciesList {
		buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: species.Name, CallbackData: strconv.Itoa(species.ID)})
	}

	return &buttons, nil
}

func (handler *Handler) breedButtons(speciesID int) (*telegram.InlineKeyboardMarkup, error) {

	breedList := make(map[int]*models.Breed)

	list, err := handler.session.GetObject(breeds)
	if err != nil {

		if errors.Is(err, ErrObjectNotExists) {
			breedList, err = handler.storage.GetBreedsList(context.Background(), speciesID)
			if err != nil {
				return nil, err
			}
			handler.session.UpdateObject(breeds, breedList)
		} else {
			return nil, err
		}

	} else {

		_, ok := list.(map[int]*models.Breed)
		if !ok {
			slog.Error("breedButtons: type assertion problem: expected breeds, get:", list)
			return nil, ErrExpectedAnotherType
		}

		breedList = list.(map[int]*models.Breed)
	}

	sortBreedList := make([]*models.Breed, 0, len(breedList))
	for _, breed := range breedList {
		sortBreedList = append(sortBreedList, breed)
	}
	sort.Slice(sortBreedList, func(i, j int) bool {
		return sortBreedList[i].Name < sortBreedList[j].Name
	})

	buttons := telegram.NewInlineKeyboardMarkup()

	for _, breed := range sortBreedList {
		buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: breed.Name, CallbackData: strconv.Itoa(breed.ID)})
	}

	return &buttons, nil
}

func (handler *Handler) sexButtons() *telegram.InlineKeyboardMarkup {

	buttons := telegram.NewInlineKeyboardMarkup()

	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "Female", CallbackData: "female"})
	buttons.AddButtonHorizontalInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "Male", CallbackData: "male"}, 0)

	return &buttons
}

func (handler *Handler) skipButton() *telegram.InlineKeyboardMarkup {

	button := telegram.NewInlineKeyboardMarkup()

	button.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "Skip", CallbackData: "skip"})

	return &button
}

func (handler *Handler) confirmationButtons() *telegram.InlineKeyboardMarkup {

	buttons := telegram.NewInlineKeyboardMarkup()

	buttons.AddButtonInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "Confirm", CallbackData: "confirm"})
	buttons.AddButtonHorizontalInlineKeyboardMarkup(&telegram.InlineKeyboardButton{Text: "Don't confirm", CallbackData: "do_not_confirm"}, 0)

	return &buttons
}
