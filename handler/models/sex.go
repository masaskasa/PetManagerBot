package models

type Sex uint

const (
	Female Sex = iota + 1
	Male
)

func (sex Sex) String() string {
	if sex == Female {
		return "⚧️ Gender: " + "♀️\n"
	} else if sex == Male {
		return "⚧️ Gender: " + "♂️\n"
	}
	return "\n-"
}
