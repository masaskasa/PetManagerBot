package sqlite

import (
	"PetManagerBot/handler"
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func NewSqliteDB(path string) (*Storage, error) {

	slog.Info("NewSqliteDB: open sqlite driven in %s", path)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		slog.Error("NewSqliteDB: can't open database:", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("NewSqliteDB: can't connect to database:", err.Error())
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (storage *Storage) Save(ctx context.Context, pet *handler.Pet) error {

	query := `insert into Pets (pet_id, owner, name, species_id, breed_id, sex) values (?,?,?,?,?,?)`

	result, err := storage.db.ExecContext(ctx, query, pet.ID, pet.Owner, pet.Name, pet.Species.ID, pet.Breed.ID, pet.Sex)
	if err != nil {
		slog.Error("Save: can't save pet:", err)
		return err
	}

	slog.Info("Save: result of sql request:", result)

	return nil
}

func (storage *Storage) IsExists(ctx context.Context, petID uuid.UUID) (bool, error) {

	query := `select exists(select 1 from Pets where pet_id = ?) as is_exist`

	var result int

	err := storage.db.QueryRowContext(ctx, query, petID).Scan(&result)
	if err != nil {
		slog.Error("IsExists: can't check pet's existence:", err)
		return false, err
	}

	return result == 1, nil
}

func (storage *Storage) Remove(ctx context.Context, petID uuid.UUID) error {

	query := `delete from Pets where pet_id = ?`

	result, err := storage.db.ExecContext(ctx, query, petID)
	if err != nil {
		slog.Error("Remove: can't remove pet:", err)
		return err
	}

	slog.Info("Remove: result of sql request:", result)

	return nil
}

func (storage *Storage) Get(ctx context.Context, petID uuid.UUID) (*handler.Pet, error) {
	//TODO implement me
	panic("implement me")
}

func (storage *Storage) Update(ctx context.Context, pet *handler.Pet) error {
	//TODO implement me
	panic("implement me")
}

func (storage *Storage) GetSpeciesList(ctx context.Context) ([]handler.Species, error) {
	//TODO implement me
	panic("implement me")
}

func (storage *Storage) GetBreedsList(ctx context.Context, speciesID int) ([]handler.Breed, error) {
	//TODO implement me
	panic("implement me")
}

func (storage *Storage) Init(ctx context.Context) error {

	query := `create table if not exists Pets (pet_id text unique not null, owner text not null, name text not null, species_id integer not null, breed_id integer not null, sex integer not null, animal_id text, special_signs text)`

	result, err := storage.db.ExecContext(ctx, query)
	if err != nil {
		slog.Error("Init: can't create table:", err)
		return err
	}

	slog.Info("Init: result of sql request:", result)

	return nil
}
