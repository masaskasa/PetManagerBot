package models

import (
	"fmt"
)

type Species struct {
	ID   int
	Name string
	Icon string
}

func (species *Species) String() string {
	return fmt.Sprintf(species.Icon+" Species: %s\n", species.Name)
}
