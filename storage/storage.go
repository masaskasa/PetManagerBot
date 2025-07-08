package storage

import (
	"PetManagerBot/handler/models"
	"context"
	"errors"
	"github.com/google/uuid"
)

type Storage interface {
	Save(ctx context.Context, pet *models.Pet) error
	IsExists(ctx context.Context, petID uuid.UUID) (bool, error)
	Remove(ctx context.Context, petID uuid.UUID) error
	Get(ctx context.Context, petID uuid.UUID) (*models.Pet, error)
	Update(ctx context.Context, pet *models.Pet) error
	GetPetsList(ctx context.Context, owner string) ([]models.Pet, error)
	GetSpeciesList(ctx context.Context) (map[int]*models.Species, error)
	GetBreedsList(ctx context.Context, speciesID int) (map[int]*models.Breed, error)
}

var ErrNoSavedPets = errors.New("no saved pets")
