package models

import "strconv"

type Breed struct {
	ID   int
	Name string
}

func (breed *Breed) String() string {
	return breed.Name + " /" + strconv.Itoa(breed.ID)
}
