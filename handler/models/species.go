package models

import "strconv"

type Species struct {
	ID     int
	Name   string
	Breeds []Breed
}

func (species *Species) String() string {
	return species.Name + " /" + strconv.Itoa(species.ID)
}

func showSpecies() []Species {
	// create list of species TODO
	return make([]Species, 0)
}

func (species *Species) showBreeds() []Breed {
	// create list of breeds TODO
	return make([]Breed, 0)
}
