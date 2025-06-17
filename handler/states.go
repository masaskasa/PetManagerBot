package handler

type scenario string

const (
	none             = "none"
	startCommand     = "/start"
	createPetCommand = "/create_pet"
	showPetCommand   = "/show_pet"
	editPetCommand   = "/edit_pet"
	deletePetCommand = "/delete_pet"
	helpCommand      = "/help"
)

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
