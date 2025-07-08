package handler

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/handler/models"
	storagePack "PetManagerBot/storage"
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func startCreatePetScenario(session *Session) {
	session.setScenario(createPetCommand)
	session.setState(start)
}

var (
	ErrUnknownSpecies = errors.New("unknown species")
	ErrUnknownBreed   = errors.New("unknown breed")
	ErrUnknownSex     = errors.New("unknown sex")
)

func createPet(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	switch session.state {
	case start:
		return setNameComplete(session, sendMessage, storage)
	case nameComplete:
		return setSpeciesComplete(session, sendMessage, storage)
	case speciesComplete:
		return setBreedComplete(session, sendMessage, storage)
	case breedComplete:
		return setSexComplete(session, sendMessage)
	case sexComplete:
		return setAnimalIDComplete(session, sendMessage)
	case animalIDComplete:
		return setSpecialSignsComplete(session, sendMessage)
	case specialSignsComplete:
		return setReadyCreatePet(session, sendMessage, storage)
	default:
		return nil
	}
}

func setNameComplete(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	owner, err := session.GetObject(userName)
	if err != nil {
		return err
	}

	newPet := models.NewPet(owner.(string))

	if err := newPet.SetName(answer.(string)); err != nil {
		return err
	}

	session.UpdateObject(pet, newPet)
	session.setState(nameComplete)

	speciesList, err := speciesList(session, storage)
	if err != nil {
		return err
	}

	_, result := sendMessage(msgAskSpecies + speciesList)
	return result
}

func setSpeciesComplete(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	species, err := determineSpecies(session)
	if errors.Is(err, ErrUnknownSpecies) {
		speciesList, err := speciesList(session, storage)
		if err != nil {
			return err
		}
		_, result := sendMessage(msgInvalidSpecies + speciesList)
		return result
	} else if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	newPet.SetSpecies(species)

	session.UpdateObject(pet, newPet)
	session.setState(speciesComplete)

	breedList, err := breedList(session, storage, species.ID)
	if err != nil {
		return err
	}

	_, result := sendMessage(msgAskBreed + breedList)
	return result
}

func setBreedComplete(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	breed, err := determineBreed(session)
	if errors.Is(err, ErrUnknownBreed) {
		breedList, err := breedList(session, storage, newPet.Species.ID)
		if err != nil {
			return err
		}
		_, result := sendMessage(msgInvalidBreed + breedList)
		return result
	} else if err != nil {
		return err
	}

	newPet.SetBreed(breed)

	session.UpdateObject(pet, newPet)
	session.setState(breedComplete)

	_, result := sendMessage(msgAskSex)
	return result
}

func setSexComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	sex, err := determineSex(session)
	if errors.Is(err, ErrUnknownSex) {
		_, result := sendMessage(msgInvalidSex)
		return result
	} else if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	newPet.SetSex(sex)

	session.UpdateObject(pet, newPet)
	session.setState(sexComplete)

	_, result := sendMessage(msgAskAnimalID)
	return result
}

func setAnimalIDComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	animalID := strings.Trim(answer.(string), "/")
	animalID = strings.ToLower(animalID)

	if animalID != "skip" {
		if err := newPet.SetAnimalID(answer.(string)); err != nil {
			return err
		}
		session.UpdateObject(pet, newPet)
	}

	session.setState(animalIDComplete)

	_, result := sendMessage(msgAskSpecialSigns)
	return result
}

func setSpecialSignsComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	specialSigns := strings.Trim(answer.(string), "/")
	specialSigns = strings.ToLower(specialSigns)

	if specialSigns != "skip" {
		newPet.SetSpecialSigns(answer.(string))
		session.UpdateObject(pet, newPet)
	}

	session.setState(specialSignsComplete)

	_, result := sendMessage(msgConfirmCreatePet + newPet.String())
	return result
}

func setReadyCreatePet(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	var result error

	confirmation := strings.Trim(answer.(string), "/")
	confirmation = strings.ToLower(confirmation)

	switch confirmation {
	case "confirm":
		if err := storage.Save(context.Background(), newPet); err != nil {
			return err
		}
		_, result = sendMessage(msgReadyCreatePet)
	case "do_not_confirm":
		_, result = sendMessage(fmt.Sprint(msgTryAgain, session.scenario))
	default:
		_, result = sendMessage(msgConfirmCreatePet + newPet.String())
		return result
	}

	session.setState(ready)
	session.setScenario(none)

	return result
}

func determinePet(session *Session) (*models.Pet, error) {

	pet, err := session.GetObject(pet)
	if err != nil {
		return nil, err
	}

	newPet := pet.(*models.Pet)

	return newPet, nil
}

func determineSpecies(session *Session) (*models.Species, error) {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return nil, err
	}

	serializedID := strings.Trim(answer.(string), "/")

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, ErrUnknownSpecies
	}

	list, err := session.GetObject(species)
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

func determineBreed(session *Session) (*models.Breed, error) {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return nil, err
	}

	serializedID := strings.Trim(answer.(string), "/")

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, ErrUnknownBreed
	}

	list, err := session.GetObject(breed)
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

func determineSex(session *Session) (models.Sex, error) {

	answer, err := session.GetObject(messageText)
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

func speciesList(session *Session, storage storagePack.Storage) (string, error) {

	speciesList := make(map[int]*models.Species)

	list, err := session.GetObject(species)
	if err != nil {
		speciesList, err = storage.GetSpeciesList(context.Background())
		if err != nil {
			return "", err
		}
		session.UpdateObject(species, speciesList)
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

func breedList(session *Session, storage storagePack.Storage, speciesID int) (string, error) {

	breedList := make(map[int]*models.Breed)

	list, err := session.GetObject(breed)
	if err != nil {
		breedList, err = storage.GetBreedsList(context.Background(), speciesID)
		if err != nil {
			return "", err
		}
		session.UpdateObject(breed, breedList)
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
