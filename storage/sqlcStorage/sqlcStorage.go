package sqlcstorage

import (
	"context"
	"database/sql"
	"errors"
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
	if err = s.Init(context.TODO()); err != nil {
		return nil, fmt.Errorf("can't initialize database: %w", err)
	}
	s.queries = db_access.New(db)
	return s, nil
}

func(s *Storage) Init(ctx context.Context) error {
	sqlInit, err := os.ReadFile("db_init.sql")
	if err != nil {
		return  fmt.Errorf("can't open sql init file: %w", err)
	}
	_, err = s.DB.ExecContext(ctx, string(sqlInit))
	if err != nil {
		return  fmt.Errorf("can't open sql init file: %w", err)
	}
	

	// q1 := `CREATE DATABASE IF NOT EXISTS phonebook;`
	// q2 := `USE phonebook;`
	// q3 := `CREATE TABLE IF NOT EXISTS person (id INT PRIMARY KEY AUTO_INCREMENT, first_name VARCHAR(50) NOT NULL,
	// last_name VARCHAR(50) NOT NULL DEFAULT 'Unknown', notes VARCHAR(200) DEFAULT 'Unknown',
	// created TIMESTAMP DEFAULT CURRENT_TIMESTAMP(), modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP());`
	// q4 := `CREATE TABLE IF NOT EXISTS phone (id INT PRIMARY KEY AUTO_INCREMENT, phone_number VARCHAR(20) UNIQUE NOT NULL, 
	// phone_type SMALLINT NOT NULL DEFAULT 3, person_id INT NOT NULL, CONSTRAINT phone_owner
	// FOREIGN KEY (person_id)	REFERENCES person(id) ON DELETE CASCADE);`

	// qs := []string{q1, q2, q3, q4}

	// for _, v := range qs {
	// 	if _, err := s.DB.ExecContext(ctx, v); err != nil {
	// 		return fmt.Errorf("can't init db: %w", err)
	// 	}
	// }
	return nil
}

func(s *Storage) Create(ctx context.Context, p *entities.Person) error {
	
	personParams := db_access.AddPersonParams{}
	personParams.FirstName = p.FirstName
	personParams.LastName = p.LastName
	personParams.Notes = sql.NullString{
		String: p.Notes,
		Valid: true,
	}
	id, err := s.queries.AddPerson(ctx, personParams)
	if err != nil {
		return fmt.Errorf("can't add person: %w", err)
	}
	phoneParams := db_access.AddPhoneParams{}
	phoneParams.PersonID = int32(id)

	for _, v := range p.Phones{
		phoneParams.PhoneNumber = v.Number
		phoneParams.PhoneType = int32(v.Type)
		err = s.queries.AddPhone(ctx, phoneParams)
		if err != nil {
			return fmt.Errorf("can't add phone: %w", err)
		}
	}
	return nil
}

func(s *Storage) Read(ctx context.Context, sd *entities.SearchData) ([]*entities.Person, error) {
	switch sd.SearchComb {
	case entities.SearchFname:
		personsDb, err := s.queries.GetPersonByFname(ctx, sd.FirstName)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoPersonExists
		} else if err != nil {
			return nil, fmt.Errorf("can't query person by fname: %w", err)
		}
		persons, err := s.convertFromModels(ctx, personsDb)
		if err != nil {
			return nil, fmt.Errorf("can't recreate person from db data: %w", err)
		}
		return persons, nil
	}
	return nil, nil
}

func(s *Storage) convertFromModels(ctx context.Context, 
	pm []db_access.GetPersonByFnameRow) ([]*entities.Person, error) {
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