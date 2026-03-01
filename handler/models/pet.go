package models

import (
	"fmt"
	"github.com/google/uuid"
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

const (
	EditName         = "EditName"
	EditSpecies      = "EditSpecies"
	EditBreed        = "EditBreed"
	EditGender       = "EditGender"
	EditAnimalID     = "EditAnimalID"
	EditSpecialSigns = "EditSpecialSigns"
)

func NewPet(owner string) *Pet {
	return &Pet{
		ID:    uuid.New(),
		Owner: owner,
	}
}

func (pet *Pet) String() string {

	name := fmt.Sprintf("🐾 *%s* 🐾\n\n", pet.Name)

	var species string
	if pet.Species != nil {
		species = fmt.Sprintf(pet.Species.Icon+" Species: %s\n", pet.Species.Name)
	}

	var breed string
	if pet.Breed != nil {
		breed = fmt.Sprintf("🏷️ Breed: %s\n", pet.Breed.Name)
	}

	var sex string
	if pet.Sex == Female {
		sex = fmt.Sprintf("⚧️ Gender: " + "♀️\n")
	} else if pet.Sex == Male {
		sex = fmt.Sprintf("⚧️ Gender: " + "♂️\n")
	}

	var animalID string
	if pet.AnimalID != "" {
		animalID = fmt.Sprintf("🆔 Animal ID: %s\n", pet.AnimalID)
	}

	var specialSigns string
	if pet.SpecialSigns != "" {
		specialSigns = fmt.Sprintf("🔍 Special Signs: %s\n", pet.SpecialSigns)
	}

	quote := "━━━━━━━━━━━━━━"

	result := name + quote + "\n" + species + breed + sex + animalID + specialSigns + quote

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

func editFlag(editingFields map[string]bool, targetParameter string) string {
	if editingFields != nil {
		for parameter, flag := range editingFields {
			switch parameter {
			case targetParameter:
				if flag {
					return "✏️"
				}
			default:
				continue
			}
		}
	}
	return ""
}
