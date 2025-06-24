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
