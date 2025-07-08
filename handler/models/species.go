package models

import "strconv"

type Species struct {
	ID   int
	Name string
}

func (species *Species) String() string {
	return species.Name + " /" + strconv.Itoa(species.ID)
}
