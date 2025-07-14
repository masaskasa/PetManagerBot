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
	ErrUnknownSpecies = errors.New("unknown species")
	ErrUnknownBreed   = errors.New("unknown breed")
	ErrUnknownSex     = errors.New("unknown sex")
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

	answer, err := handler.session.GetObject(messageText)
	if err != nil {
		return err
	}

	owner, err := handler.session.GetObject(userName)
	if err != nil {
		return err
	}

	newPet := models.NewPet(owner.(string))

	if err := newPet.SetName(answer.(string)); err != nil {
		return err
	}

	handler.session.UpdateObject(pet, newPet)
	handler.session.setState(nameComplete)

	speciesButtons, err := handler.speciesButtons()
	if err != nil {
		return err
	}

	_, result := handler.sendMessageKeyboard(msgAskSpecies, *speciesButtons)
	return result
}

func (handler *Handler) setSpeciesComplete() error {

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

	newPet, err := handler.determinePet()
	if err != nil {
		return err
	}

	newPet.SetSpecies(species)

	if handler.answerCallbackQuery != nil {
		if _, err := handler.answerCallbackQuery(ntfSetSpecies+species.Name, false); err != nil {
			slog.Error("setSpeciesComplete: answerCallbackQuery:", err)
		}
	}

	handler.session.UpdateObject(pet, newPet)
	handler.session.setState(speciesComplete)

	breedList, err := handler.breedList(species.ID)
	if err != nil {
		return err
	}

	_, result := handler.sendMessage(msgAskBreed + breedList)
	return result
}

func (handler *Handler) setBreedComplete() error {

	newPet, err := handler.determinePet()
	if err != nil {
		return err
	}

	breed, err := handler.determineBreed()
	if errors.Is(err, ErrUnknownBreed) {
		breedList, err := handler.breedList(newPet.Species.ID)
		if err != nil {
			return err
		}
		_, result := handler.sendMessage(msgInvalidBreed + breedList)
		return result
	} else if err != nil {
		return err
	}

	newPet.SetBreed(breed)

	handler.session.UpdateObject(pet, newPet)
	handler.session.setState(breedComplete)

	_, result := handler.sendMessage(msgAskSex)
	return result
}

func (handler *Handler) setSexComplete() error {

	sex, err := handler.determineSex()
	if errors.Is(err, ErrUnknownSex) {
		_, result := handler.sendMessage(msgInvalidSex)
		return result
	} else if err != nil {
		return err
	}

	newPet, err := handler.determinePet()
	if err != nil {
		return err
	}

	newPet.SetSex(sex)

	handler.session.UpdateObject(pet, newPet)
	handler.session.setState(sexComplete)

	_, result := handler.sendMessage(msgAskAnimalID)
	return result
}

func (handler *Handler) setAnimalIDComplete() error {

	answer, err := handler.session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := handler.determinePet()
	if err != nil {
		return err
	}

	animalID := strings.Trim(answer.(string), "/")
	animalID = strings.ToLower(animalID)

	if animalID != "skip" {
		if err := newPet.SetAnimalID(answer.(string)); err != nil {
			return err
		}
		handler.session.UpdateObject(pet, newPet)
	}

	handler.session.setState(animalIDComplete)

	_, result := handler.sendMessage(msgAskSpecialSigns)
	return result
}

func (handler *Handler) setSpecialSignsComplete() error {

	answer, err := handler.session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := handler.determinePet()
	if err != nil {
		return err
	}

	specialSigns := strings.Trim(answer.(string), "/")
	specialSigns = strings.ToLower(specialSigns)

	if specialSigns != "skip" {
		newPet.SetSpecialSigns(answer.(string))
		handler.session.UpdateObject(pet, newPet)
	}

	handler.session.setState(specialSignsComplete)

	_, result := handler.sendMessage(msgConfirmCreatePet + newPet.String())
	return result
}

func (handler *Handler) setReadyCreatePet() error {

	answer, err := handler.session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := handler.determinePet()
	if err != nil {
		return err
	}

	var result error

	confirmation := strings.Trim(answer.(string), "/")
	confirmation = strings.ToLower(confirmation)

	switch confirmation {
	case "confirm":
		if err := handler.storage.Save(context.Background(), newPet); err != nil {
			return err
		}
		_, result = handler.sendMessage(msgReadyCreatePet)
	case "do_not_confirm":
		_, result = handler.sendMessage(fmt.Sprint(msgTryAgain, handler.session.scenario))
	default:
		_, result = handler.sendMessage(msgConfirmCreatePet + newPet.String())
		return result
	}

	handler.session.setState(ready)
	handler.session.setScenario(none)

	return result
}

func (handler *Handler) determinePet() (*models.Pet, error) {

	pet, err := handler.session.GetObject(pet)
	if err != nil {
		return nil, err
	}

	newPet := pet.(*models.Pet)

	return newPet, nil
}

func (handler *Handler) determineSpecies() (*models.Species, error) {

	var serializedID string

	answer, err := handler.session.GetObject(messageText)
	if err == nil {
		serializedID = strings.Trim(answer.(string), "/")
	}

	answer, err = handler.session.GetObject(callbackQueryData)
	if err == nil {
		serializedID = answer.(string)
	}

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, ErrUnknownSpecies
	}

	list, err := handler.session.GetObject(species)
	if err != nil {
		return nil, err
	}

	speciesList := list.(map[int]*models.Species)

	species, exists := speciesList[id]
	if !exists {
		return nil, ErrUnknownSpecies
	}

	return species, nil
}

func (handler *Handler) determineBreed() (*models.Breed, error) {

	answer, err := handler.session.GetObject(messageText)
	if err != nil {
		return nil, err
	}

	serializedID := strings.Trim(answer.(string), "/")

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, ErrUnknownBreed
	}

	list, err := handler.session.GetObject(breed)
	if err != nil {
		return nil, err
	}

	breedList := list.(map[int]*models.Breed)

	breed, exists := breedList[id]
	if !exists {
		return nil, ErrUnknownBreed
	}

	return breed, nil
}

func (handler *Handler) determineSex() (models.Sex, error) {

	answer, err := handler.session.GetObject(messageText)
	if err != nil {
		return 0, err
	}

	sex := strings.Trim(answer.(string), "/")
	sex = strings.ToLower(sex)

	switch sex {
	case "female":
		return models.Female, nil
	case "male":
		return models.Male, nil
	default:
		return 0, ErrUnknownSex
	}
}

func (handler *Handler) speciesList() (string, error) {

	speciesList := make(map[int]*models.Species)

	list, err := handler.session.GetObject(species)
	if err != nil {
		speciesList, err = handler.storage.GetSpeciesList(context.Background())
		if err != nil {
			return "", err
		}
		handler.session.UpdateObject(species, speciesList)
	} else {
		speciesList = list.(map[int]*models.Species)
	}

	words := make([]string, 0, 10)

	for _, species := range speciesList {
		words = append(words, species.String())
	}

	serializedList := concatBuilder(words)

	return serializedList, nil
}

func (handler *Handler) breedList(speciesID int) (string, error) {

	breedList := make(map[int]*models.Breed)

	list, err := handler.session.GetObject(breed)
	if err != nil {
		breedList, err = handler.storage.GetBreedsList(context.Background(), speciesID)
		if err != nil {
			return "", err
		}
		handler.session.UpdateObject(breed, breedList)
	} else {
		breedList = list.(map[int]*models.Breed)
	}

	words := make([]string, 0, 20)

	for _, breed := range breedList {
		words = append(words, breed.String())
	}

	serializedList := concatBuilder(words)

	return serializedList, nil
}

func concatBuilder(words []string) string {

	sort.Strings(words)

	var b strings.Builder
	for _, word := range words {
		b.WriteString(word + "\n")
	}

	return b.String()
}

func (handler *Handler) speciesButtons() (*telegram.InlineKeyboardMarkup, error) {

	speciesList := make(map[int]*models.Species)

	list, err := handler.session.GetObject(species)
	if err != nil {
		speciesList, err = handler.storage.GetSpeciesList(context.Background())
		if err != nil {
			return nil, err
		}
		handler.session.UpdateObject(species, speciesList)
	} else {
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
