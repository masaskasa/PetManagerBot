package handler

import (
	"PetManagerBot/clients/telegram"
	"PetManagerBot/handler/models"
	storagePack "PetManagerBot/storage"
	"context"
	"fmt"
	"strconv"
	"strings"
)

func doCreatePetScenario(session *Session) {
	session.setScenario(createPetCommand)
	session.setState(start)
}

func createPet(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	switch session.state {
	case start:
		return doNameCompleteState(session, sendMessage, storage)
	case nameComplete:
		return doSpeciesComplete(session, sendMessage, storage)
	case speciesComplete:
		return doBreedComplete(session, sendMessage)
	case breedComplete:
		return doSexComplete(session, sendMessage)
	case sexComplete:
		return doAnimalIDComplete(session, sendMessage)
	case animalIDComplete:
		return doSpecialSignsComplete(session, sendMessage)
	case specialSignsComplete:
		return doReadyCreatePet(session, sendMessage, storage)
	default:
		return nil
	}
}

func doNameCompleteState(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	name, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	owner, err := session.GetObject(userName)
	if err != nil {
		return err
	}

	newPet := models.NewPet(owner.(string))

	if err := newPet.SetName(name.(string)); err != nil {
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

func doSpeciesComplete(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	species, err := determineSpecies(session)
	if err != nil {
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

func doBreedComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	breed, err := determineBreed(session)
	if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	newPet.SetBreed(breed)

	session.UpdateObject(pet, newPet)

	session.setState(breedComplete)

	_, result := sendMessage(msgAskSex)
	return result
}

func doSexComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	gender, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	var sex models.Sex

	if gender == "/female" {
		sex = models.Female
	} else if gender == "/male" {
		sex = models.Male
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

func doAnimalIDComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	if answer != "/skip" {

		animalID := answer.(string)

		newPet, err := determinePet(session)
		if err != nil {
			return err
		}

		if err := newPet.SetAnimalID(animalID); err != nil {
			return err
		}

		session.UpdateObject(pet, newPet)
	}

	session.setState(animalIDComplete)

	_, result := sendMessage(msgAskSpecialSigns)
	return result
}

func doSpecialSignsComplete(session *Session, sendMessage func(string) (telegram.Message, error)) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	if answer != "/skip" {

		specialSigns := answer.(string)

		newPet.SetSpecialSigns(specialSigns)

		session.UpdateObject(pet, newPet)
	}

	session.setState(specialSignsComplete)

	_, result := sendMessage(msgConfirmCreatePet + newPet.String())
	return result
}

func doReadyCreatePet(session *Session, sendMessage func(string) (telegram.Message, error), storage storagePack.Storage) error {

	answer, err := session.GetObject(messageText)
	if err != nil {
		return err
	}

	newPet, err := determinePet(session)
	if err != nil {
		return err
	}

	var result error

	if answer == "/confirm" {

		if err := storage.Save(context.Background(), newPet); err != nil {
			return err
		}

		_, result = sendMessage(msgReadyCreatePet)
	} else {
		_, result = sendMessage(fmt.Sprint(msgTryAgain, session.scenario))
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

	code, err := session.GetObject(messageText)
	if err != nil {
		return nil, err
	}

	serializedID := strings.Trim(code.(string), "/")

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, err
	}

	list, err := session.GetObject(species)
	if err != nil {
		return nil, err
	}

	speciesList := list.([]models.Species)

	for _, species := range speciesList {
		if species.ID == id {
			return &species, nil
		}
	}
	return nil, nil
}

func determineBreed(session *Session) (*models.Breed, error) {

	code, err := session.GetObject(messageText)
	if err != nil {
		return nil, err
	}

	serializedID := strings.Trim(code.(string), "/")

	id, err := strconv.Atoi(serializedID)
	if err != nil {
		return nil, err
	}

	list, err := session.GetObject(breed)
	if err != nil {
		return nil, err
	}

	speciesList := list.([]models.Breed)

	for _, breed := range speciesList {
		if breed.ID == id {
			return &breed, nil
		}
	}
	return nil, nil
}

func speciesList(session *Session, storage storagePack.Storage) (string, error) {

	list, err := storage.GetSpeciesList(context.Background())
	if err != nil {
		return "", err
	}

	session.UpdateObject(species, list)

	words := make([]string, 0, 10)

	for _, species := range list {
		words = append(words, species.String())
	}

	concatBuilder(words)
	serializedList := concatBuilder(words)

	return serializedList, nil
}

func concatBuilder(words []string) string {

	var b strings.Builder
	for _, word := range words {
		b.WriteString(word + "\n")
	}

	return b.String()
}

func breedList(session *Session, storage storagePack.Storage, speciesID int) (string, error) {

	list, err := storage.GetBreedsList(context.Background(), speciesID)
	if err != nil {
		return "", err
	}

	session.UpdateObject(breed, list)

	words := make([]string, 0, 20)

	for _, breed := range list {
		words = append(words, breed.String())
	}

	concatBuilder(words)
	serializedList := concatBuilder(words)

	return serializedList, nil
}
