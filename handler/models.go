package handler

import (
	"github.com/google/uuid"
	"log/slog"
)

type Pet struct {
	ID      uuid.UUID
	Owner   string
	Name    string
	Species *Species
	Breed   *Breed
	Sex     Sex
	// Photo TODO
	AnimalID     string
	SpecialSigns string
}

func NewPet(owner string) *Pet {
	return &Pet{
		ID:    uuid.New(),
		Owner: owner,
	}
}

func (pet *Pet) String() string {

	var species string
	if pet.Species != nil {
		species = "\n" + pet.Species.Name
	}

	var breed string
	if pet.Breed != nil {
		breed = "\n" + pet.Breed.Name
	}

	var sex string
	if pet.Sex == Female {
		sex = "\nДевочка"
	} else if pet.Sex == Male {
		sex = "\nМальчик"
	}

	var animalID string
	if pet.AnimalID != "" {
		animalID = "\n" + pet.AnimalID
	}

	var specialSigns string
	if pet.SpecialSigns != "" {
		specialSigns = "\n\n" + pet.SpecialSigns
	}

	result := pet.Name + species + breed + sex + animalID + specialSigns

	slog.Info("Pet: String() result:", result)
	return result
}

func (pet *Pet) SetName(name string) error {
	// validate TODO
	pet.Name = name
	return nil
}

func (pet *Pet) SetSpecies(species *Species) {
	pet.Species = species
}

func (pet *Pet) SetBreed(breed *Breed) {
	pet.Breed = breed
}

func (pet *Pet) SetSex(sex Sex) {
	pet.Sex = sex
}

func (pet *Pet) SetAnimalID(animalID string) error {
	// validate TODO
	pet.AnimalID = animalID
	return nil
}

func (pet *Pet) SetSpecialSigns(specialSigns string) {
	pet.SpecialSigns = specialSigns
}

type Species struct {
	ID     int
	Name   string
	Breeds []Breed
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
	ID   int
	Name string
}

type Sex uint

const (
	None Sex = iota + 1
	Female
	Male
)
