package storage

import (
	"PetManagerBot/handler"
	"context"
	"errors"
	"github.com/google/uuid"
)

type Storage interface {
	Save(ctx context.Context, pet *handler.Pet) error
	IsExists(ctx context.Context, petID uuid.UUID) (bool, error)
	Remove(ctx context.Context, petID uuid.UUID) error
	Get(ctx context.Context, petID uuid.UUID) (*handler.Pet, error)
	Update(ctx context.Context, pet *handler.Pet) error
	GetPetsList(ctx context.Context, owner string) ([]handler.Pet, error)
	GetSpeciesList(ctx context.Context) ([]handler.Species, error)
	GetBreedsList(ctx context.Context, speciesID int) ([]handler.Breed, error)
}

var ErrNoSavedPets = errors.New("no saved pets")
