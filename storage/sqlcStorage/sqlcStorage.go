package sqlcstorage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
	"github.com/bondar-aleksandr/phonebook/db_access"
	"github.com/bondar-aleksandr/phonebook/entities"
	"github.com/bondar-aleksandr/phonebook/storage"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	DB *sql.DB
	queries *db_access.Queries
}

// New creates new instance of storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("mysql", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	s := &Storage{DB: db}
	if err = s.init(context.TODO()); err != nil {
		return nil, fmt.Errorf("can't initialize database: %w", err)
	}
	s.queries = db_access.New(db)
	return s, nil
}

// Init executes commands from sql_init.sql file. Creates database
// and tables.
func(s *Storage) init(ctx context.Context) error {
	sqlInit, err := os.ReadFile("sql/db_init.sql")
	if err != nil {
		return  fmt.Errorf("can't open sql init file: %w", err)
	}
	_, err = s.DB.ExecContext(ctx, string(sqlInit))
	if err != nil {
		return  fmt.Errorf("can't execute sql init file: %w", err)
	}
	return nil
}

// Reset executes commands from sql_reset.sql file. Drops database.
func(s *Storage) Reset(ctx context.Context) error {
	sqlreset, err := os.ReadFile("sql/db_reset.sql")
	if err != nil {
		return  fmt.Errorf("can't open sql reset file: %w", err)
	}
	_, err = s.DB.ExecContext(ctx, string(sqlreset))
	if err != nil {
		return  fmt.Errorf("can't execute sql reset file: %w", err)
	}
	return nil
}

// Create CRUD operation.
func(s *Storage) Create(ctx context.Context, p *entities.Person) error {
	
	personParams := db_access.AddPersonParams{}
	personParams.FirstName = p.FirstName
	personParams.LastName = p.LastName
	if p.Notes != "" {
		personParams.Notes = sql.NullString{
			String: p.Notes,
			Valid: true,
		}
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	id, err := qtx.AddPerson(ctx, personParams)
	if err != nil {
		return fmt.Errorf("can't add person: %w", err)
	}
	phoneParams := db_access.AddPhoneParams{}
	phoneParams.PersonID = int32(id)

	for _, v := range p.Phones{
		phoneParams.PhoneNumber = v.Number
		phoneParams.PhoneType = int32(v.Type)
		err = qtx.AddPhone(ctx, phoneParams)
		if err != nil {
			return fmt.Errorf("can't add phone: %w", err)
		}
	}
	return tx.Commit()
}

// Read CRUD operation.
// TODO: refactor if possible, don't to repeat the same actions for every case!
func(s *Storage) Read(ctx context.Context, cd *storage.CrudData) ([]*entities.Person, error) {
	var persons []*entities.Person
	switch cd.SearchComb {
	case storage.CrudFname:
		personsDb, err := s.queries.GetPersonByFname(ctx, cd.FirstName)
		if err != nil {
			return nil, fmt.Errorf("can't query person by fname: %w", err)
		}
		persons, err = s.convertFromDBModels(ctx, personsDb)
		if err != nil {
			return nil, fmt.Errorf("can't recreate person from db data: %w", err)
		}
	case storage.CrudLname:
		personsDb, err := s.queries.GetPersonByLname(ctx, cd.LastName)
		if err != nil {
			return nil, fmt.Errorf("can't query person by lname: %w", err)
		}
		persons, err = s.convertFromDBModels(ctx, personsDb)
		if err != nil {
			return nil, fmt.Errorf("can't recreate person from db data: %w", err)
		}
	case storage.CrudAll:
		personsDb, err := s.queries.GetAllPersons(ctx)
		if err != nil {
			return nil, fmt.Errorf("can't query persons: %w", err)
		}
		persons, err = s.convertFromDBModels(ctx, personsDb)
		if err != nil {
			return nil, fmt.Errorf("can't recreate person from db data: %w", err)
		}
	case storage.CrudPhone:
		ids, err := s.queries.GetPersonIDByPhone(ctx, cd.Phone)
		if err != nil {
			return nil, fmt.Errorf("can't query persons id: %w", err)
		}
		personsDb := []db_access.Person{}
		for _, id := range ids {
			person, err := s.queries.GetPersonByID(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("can't query person by id: %w", err)
			}
			personsDb = append(personsDb, person)
		}
		persons, err = s.convertFromDBModels(ctx, personsDb)
		if err != nil {
			return nil, fmt.Errorf("can't recreate person from db data: %w", err)
		}
	}
	return persons, nil
}

// convertFromDBModel queries phone info for person, and recreates *entities.Person
// objects 
func(s *Storage) convertFromDBModels(ctx context.Context, 
	pm []db_access.Person) ([]*entities.Person, error) {
	result := []*entities.Person{}
	for _,v := range pm {
		p := entities.NewPerson(v.FirstName, v.LastName, v.Notes.String)
		phones, err := s.queries.GetPersonPhones(ctx, v.ID)
		if err != nil {
			return nil, fmt.Errorf("can't query person phone: %w", err)
		}
		for _,phone := range phones {
			ph := entities.NewPhone(entities.PhoneType(phone.PhoneType), phone.PhoneNumber)
			p.AddPhone(ph)
		}
		result = append(result, p)
	}
	return result, nil
}