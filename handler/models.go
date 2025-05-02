package handler

import (
	"github.com/google/uuid"
)

type Pet struct {
	id      uuid.UUID
	owner   string
	name    string
	species *Species
	breed   *Breed
	sex     Sex
	// photo TODO
	animalID     string
	specialSigns string
}

func newPet(owner string) *Pet {
	return &Pet{
		id:    uuid.New(),
		owner: owner,
	}
}

func (pet *Pet) setName(name string) error {
	// validate TODO
	pet.name = name
	return nil
}

func (pet *Pet) setSpecies(species *Species) {
	pet.species = species
}

func (pet *Pet) setBreed(breed *Breed) {
	pet.breed = breed
}

func (pet *Pet) setSex(sex Sex) {
	pet.sex = sex
}

func (pet *Pet) setAnimalID(animalID string) error {
	// validate TODO
	pet.animalID = animalID
	return nil
}

func (pet *Pet) setSpecialSigns(specialSigns string) {
	pet.specialSigns = specialSigns
}

type Species struct {
	id     uuid.UUID
	name   string
	breeds []Breed
}

func showSpecies() []Species {
	// create list of species TODO
	return make([]Species, 0)
}

func (species *Species) showBreeds() []Breed {
	// create list of breeds TODO
	return make([]Breed, 0)
}

type Breed struct {
	id   uuid.UUID
	name string
}

type Sex uint

const (
	female Sex = iota + 1
	male
)
