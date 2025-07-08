package models

import "strconv"

type Breed struct {
	ID        int
	Name      string
	SpeciesID int
}

func (breed *Breed) String() string {
	return breed.Name + " /" + strconv.Itoa(breed.ID)
}
