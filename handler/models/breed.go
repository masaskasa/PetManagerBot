package models

import (
	"fmt"
)

type Breed struct {
	ID        int
	Name      string
	SpeciesID int
}

func (breed *Breed) String() string {
	return fmt.Sprintf("🏷️ Breed: %s\n", breed.Name)
}
