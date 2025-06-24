package models

type Sex uint

const (
	Female Sex = iota + 1
	Male
)

func (sex Sex) String() string {
	if sex == Female {
		return "\nДевочка"
	} else if sex == Male {
		return "\nМальчик"
	}
	return "\n-"
}
