package handler

const (
	msgHello                = "Hi there!👋\n\n" + msgHowToBegin
	msgAskName              = "What is your pet's name?"
	msgAskSpecies           = "What is your pet's species? Choose the correct one\n\n"
	msgInvalidSpecies       = "Invalid species. Please choose the correct one\n\n"
	ntfSetSpecies           = "You have chosen the species of "
	msgAskBreed             = "What is your pet's breed? Choose the correct one\n\n"
	msgInvalidBreed         = "Invalid breed. Please choose the correct one\n\n"
	ntfSetBreed             = "You have chosen the breed of "
	msgAskSex               = "What is your pet's gender?"
	msgInvalidSex           = "Invalid gender. Please choose the correct one"
	ntfSetSex               = "You have chosen the gender of "
	msgAskAnimalID          = "What is your pet's AnimalID?"
	msgAskSpecialSigns      = "Specify the special characteristics of the pet"
	msgConfirmCreatePet     = "Your pet is ready! if everything is correct, confirm the creation of a pet\n\n"
	msgReadyCreatePet       = "Your pet is saved! ✅"
	msgTryAgain             = "OK! You can try again "
	msgUnknownCommand       = "Unknown command 🚫"
	msgBreakCommand         = "Operation cancelled\n\n"
	msgNeedlessBreakCommand = "You are not performing any operation to cancel it"
	ntfSkip                 = "You skipped this step. You can complete it later"
	msgNoSavedPets          = "You have no saved pets ⚠️\n\nYou can create pet with command /create_pet"
	msgPickPet              = "Please select a pet"
	msgPickParameter        = "Please select a parameter for edit %s's card:\n\n"
	msgInvalidPet           = "Invalid pet. Please choose the correct one\n\n"
	msgConfirmDeletePet     = "Are you sure you want to delete the pet card? This action cannot be undone\n\n"
	msgReadyDeletePet       = "Pet's card has been deleted"
	msgCancelDelete         = "Pet's card deletion canceled"
	msgInvalidParameter     = "Invalid parameter. Please choose the correct one\n\n"

	msgHowToBegin = `How to begin?

- Create an animal card /create_pet`
	msgHelp = `Hi! I’m your pet assistant bot 🐾
Here’s what I can do:

Pet management:
/create_pet — Create a new pet profile
/show_pet — Show a pet profile (or list all pets)
/edit_pet — Edit a pet profile
/delete_pet — Delete a pet profile

Other:
/break — Cancel the current action
/help — Show this help message
`
)
