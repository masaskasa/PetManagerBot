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

	query := `select P.owner,
					 P.name,
					 P.species_id,
					 S.name,
					 P.breed_id,
					 B.name,
					 P.sex,
					 P.animal_id,
					 P.special_signs
			  from Pets as P
			  join Species as S on P.species_id=S.species_id
			  join Breeds as B on P.breed_id=B.breed_id
			  where pet_id = ?`

	var owner, name, speciesName, breedName string
	var animalIDBytes, specialSignsBytes []byte
	var speciesID, breedID, sex int

	err := storage.db.QueryRowContext(ctx, query, petID).Scan(&owner, &name, &speciesID, &speciesName, &breedID, &breedName, &sex, &animalIDBytes, &specialSignsBytes)
	if err != nil {
		slog.Error("Get: can't get pet:", err)
		return nil, err
	}

	var specialSigns string
	if specialSignsBytes != nil {
		specialSigns = string(specialSignsBytes)
	}

	var animalID string
	if animalIDBytes != nil {
		specialSigns = string(animalIDBytes)
	}

	return &handler.Pet{
		ID:    petID,
		Owner: owner,
		Name:  name,
		Species: &handler.Species{
			ID:   speciesID,
			Name: speciesName},
		Breed: &handler.Breed{
			ID:   breedID,
			Name: breedName},
		Sex:          handler.Sex(sex),
		SpecialSigns: specialSigns,
		AnimalID:     animalID,
	}, nil
}

func (storage *Storage) Update(ctx context.Context, pet *handler.Pet) error {
	//TODO implement me
	panic("implement me")
}

func (storage *Storage) GetPetsList(ctx context.Context, owner string) ([]handler.Pet, error) {

	pets := make([]handler.Pet, 0, 10)

	query := `select pet_id, name from Pets where owner = ?`

	rows, err := storage.db.Query(query, owner)
	defer func() { _ = rows.Close() }()
	if err != nil {
		slog.Error("GetPetsList: can't get pets list:", err)
		return nil, err
	}

	for rows.Next() {
		var id string
		var name string
		var petID uuid.UUID

		err = rows.Scan(&id, &name)
		if err != nil {
			slog.Error("GetPetsList: can't parse pet from row:", err)
			return nil, err
		}

		petID, err = uuid.Parse(id)
		if err != nil {
			slog.Error("GetPetsList: can't parse pet's uuid:", err)
			return nil, err
		}

		pets = append(pets, handler.Pet{ID: petID, Name: name})
	}

	slog.Info("GetPetsList: result:", pets)
	return pets, nil
}

func (storage *Storage) GetSpeciesList(ctx context.Context) ([]handler.Species, error) {

	species := make([]handler.Species, 0, 15)

	query := `select species_id, name from Species`

	rows, err := storage.db.Query(query)
	defer func() { _ = rows.Close() }()
	if err != nil {
		slog.Error("GetSpeciesList: can't get species list:", err)
		return nil, err
	}

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			slog.Error("GetSpeciesList: can't parse species from row:", err)
			return nil, err
		}

		species = append(species, handler.Species{ID: id, Name: name})
	}

	slog.Info("GetSpeciesList: result:", species)
	return species, nil
}

func (storage *Storage) GetBreedsList(ctx context.Context, speciesID int) ([]handler.Breed, error) {

	breeds := make([]handler.Breed, 0, 25)

	query := `select breed_id, name from Breeds where species_id = ?`

	rows, err := storage.db.Query(query, speciesID)
	defer func() { _ = rows.Close() }()
	if err != nil {
		slog.Error("GetBreedsList: can't get breeds list:", err)
		return nil, err
	}

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			slog.Error("GetBreedsList: can't parse breed from row:", err)
			return nil, err
		}

		breeds = append(breeds, handler.Breed{ID: id, Name: name})
	}

	slog.Info("GetBreedsList: result:", breeds)
	return breeds, nil
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
