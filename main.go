package main

import (
	"PetManagerBot/handler"
	"PetManagerBot/storage/sqlite"
	"context"
	"flag"
	"fmt"
	"log"
)

const tgHost = "api.telegram.org"

func main() {
	//client := telegram.NewClient(tgHost, mustToken())

	storage, err := sqlite.NewSqliteDB("/home/user/storage_pet_manager_bot/storage.db")
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: %s", err)
	}

	//offset := 0
	//var chatID int
	//var user string

	//for {
	//	updates, _ := client.GetUpdates(offset, 100)
	//	if len(updates) != 0 {
	//		offset = updates[len(updates)-1].ID + 1
	//		for _, upd := range updates {
	//			slog.Debug(upd.Message.Text)
	//		}
	//		chatID = updates[0].Message.Chat.ID
	//		//user = updates[0].Message.From.UserName
	//		time.Sleep(1 * time.Second)
	//	} else {
	//		continue
	//	}
	//
	//	//receivedMessage, err := client.SendMessage(chatID, "hello handsome")
	//	//if err != nil {
	//	//	log.Fatal("error send message")
	//	//}
	//	//println(receivedMessage.Text)
	//
	//	//myPet := handler.NewPet(user)
	//	//_ = myPet.SetName("Котя")
	//	//myPet.SetSpecies(&handler.Species{ID: 2,
	//	//	Name:   "Кошка",
	//	//	Breeds: make([]handler.Breed, 0)})
	//	//myPet.SetBreed(&handler.Breed{ID: 11,
	//	//	Name: "Метис"})
	//	//myPet.SetSex(handler.Female)
	//	//
	//	//if err := storage.Save(context.TODO(), myPet); err != nil {
	//	//	slog.Debug("can't save pet:", err)
	//	//}
	//
	//	result, err := storage.GetSpeciesList(context.TODO())
	//	if err != nil {
	//		log.Fatal("error get pets")
	//	}
	//
	//	var builder strings.Builder
	//	for i := range result {
	//		builder.WriteString(result[i].String() + "\n")
	//	}
	//
	//	receivedMessage, err := client.SendMessage(chatID, builder.String())
	//	if err != nil {
	//		log.Fatal("error send message")
	//	}
	//	println(receivedMessage.Text)
	//
	//	builder.Reset()
	//
	//	var resultBreeds []handler.Breed
	//	resultBreeds, err = storage.GetBreedsList(context.TODO(), 2)
	//	if err != nil {
	//		log.Fatal("error get pets")
	//	}
	//
	//	for i := range resultBreeds {
	//		builder.WriteString(resultBreeds[i].String() + "\n")
	//	}
	//
	//	receivedMessage, err = client.SendMessage(chatID, builder.String())
	//	if err != nil {
	//		log.Fatal("error send message")
	//	}
	//	println(receivedMessage.Text)
	//
	//	builder.Reset()
	//
	//	var resultPets []handler.Pet
	//	resultPets, err = storage.GetPetsList(context.TODO(), "cinamorollya")
	//	if err != nil {
	//		log.Fatal("error get pets")
	//	}
	//
	//	for i := range resultPets {
	//		builder.WriteString(resultPets[i].Name + " /" + resultPets[i].ID.String() + "\n")
	//	}
	//
	//	receivedMessage, err = client.SendMessage(chatID, builder.String())
	//	if err != nil {
	//		log.Fatal("error send message")
	//	}
	//	println(receivedMessage.Text)
	//}

	var resultPets []handler.Pet
	resultPets, err = storage.GetPetsList(context.TODO(), "aboba")
	if err != nil {
		log.Fatal("error get pets")
	}

	fmt.Println(resultPets)

	//result, err := storage.IsExists(context.TODO(), petID)
	//if err != nil {
	//	log.Fatalf("IsExists don't work: %s", err)
	//}
	//if result {
	//	receivedMessage, _ := client.SendMessage(chatID, "Котя exists!")
	//	println(receivedMessage.Text)
	//}
	//
	//if err := storage.Remove(context.TODO(), petID); err != nil {
	//	log.Fatalf("Remove don't work: %s", err)
	//}
	//if result {
	//	receivedMessage, _ := client.SendMessage(chatID, "Котя удален")
	//	println(receivedMessage.Text)
	//}

	//petID, _ := uuid.Parse("b3f6a74e-3ab7-448f-8836-1dad588dff4b")
	//
	//newPet := handler.Pet{ID: petID,
	//	Owner:    "Masha",
	//	Sex:      handler.Male,
	//	Name:     "Енот",
	//	AnimalID: "test",
	//}
	//
	//err = storage.Update(context.TODO(), &newPet)
}

/*
Return a token for client.
Token parse from command line argument with flag -token-for-bot.
If argument empty print error and call to os.Exit(1)
*/
func mustToken() string {
	token := flag.String("token-for-bot", "", "set token for telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
