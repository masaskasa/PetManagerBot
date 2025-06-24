package handler

type dialogState uint

const (
	start = iota + 1
	nameComplete
	speciesComplete
	breedComplete
	sexComplete
	animalIDComplete
	specialSignsComplete
	ready
	pickPet
	pickParameter
	parameterComplete
)
