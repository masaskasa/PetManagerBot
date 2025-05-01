package handler

import (
	"github.com/google/uuid"
)

type pet struct {
	id      uuid.UUID
	owner   string
	name    string
	species *species
	breed   *breed
	sex     sex
	// photo TODO
	animalID     string
	specialSigns string
}

func newPet(owner string) *pet {
	return &pet{
		id:    uuid.New(),
		owner: owner,
	}
}

func (pet *pet) setName(name string) error {
	// validate TODO
	pet.name = name
	return nil
}

func (pet *pet) setSpecies(species *species) {
	pet.species = species
}

func (pet *pet) setBreed(breed *breed) {
	pet.breed = breed
}

func (pet *pet) setSex(sex sex) {
	pet.sex = sex
}

func (pet *pet) setAnimalID(animalID string) error {
	// validate TODO
	pet.animalID = animalID
	return nil
}

func (pet *pet) setSpecialSigns(specialSigns string) {
	pet.specialSigns = specialSigns
}

type species struct {
	id     uuid.UUID
	name   string
	breeds []breed
}

func showSpecies() []species {
	// create list of species TODO
	return make([]species, 0)
}

func (species *species) showBreeds() []breed {
	// create list of breeds TODO
	return make([]breed, 0)
}

type breed struct {
	id   uuid.UUID
	name string
}

type sex uint

const (
	female sex = iota + 1
	male
)
