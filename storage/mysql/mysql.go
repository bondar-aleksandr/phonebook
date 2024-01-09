package mysql

import (
	"context"
	"database/sql"
	// "errors"
	"fmt"
	// "log"
	"time"

	"github.com/bondar-aleksandr/phonebook/entities"
	"github.com/bondar-aleksandr/phonebook/storage"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	DB *sql.DB
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
	if err = s.Init(context.TODO()); err != nil {
		return nil, fmt.Errorf("can't initialize database: %w", err)
	}
	return s, nil
}

func(s *Storage) Init(ctx context.Context) error {
	q1 := `CREATE DATABASE IF NOT EXISTS phonebook;`
	q2 := `USE phonebook;`
	q3 := `CREATE TABLE IF NOT EXISTS person (id INT PRIMARY KEY AUTO_INCREMENT, first_name VARCHAR(50) NOT NULL,
	last_name VARCHAR(50) NOT NULL DEFAULT 'Unknown', notes VARCHAR(200) DEFAULT 'Unknown',
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP(), modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(),
	CONSTRAINT unique_person UNIQUE(first_name, last_name));`
	q4 := `CREATE TABLE IF NOT EXISTS phone (id INT PRIMARY KEY AUTO_INCREMENT, number VARCHAR(20) UNIQUE NOT NULL, 
	type SMALLINT NOT NULL DEFAULT 3, active BOOL NOT NULL DEFAULT true, person_id INT, CONSTRAINT phone_owner
	FOREIGN KEY (person_id)	REFERENCES person(id) ON DELETE CASCADE);`

	qs := []string{q1, q2, q3, q4}

	for _, v := range qs {
		if _, err := s.DB.ExecContext(ctx, v); err != nil {
			return fmt.Errorf("can't init db: %w", err)
		}
	}
	return nil
}

func(s *Storage) Create(ctx context.Context, p *entities.Person) error {
	q := `INSERT INTO person (first_name, last_name, notes) VALUES (?, ?, ?)`

	res, err := s.DB.ExecContext(ctx, q, p.FirstName, p.LastName, p.Notes); 
	if err != nil {
		return fmt.Errorf("can't create person: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("can't get last insert id: %w", err)
	}
	stmt, err := s.DB.Prepare(`INSERT INTO phone (number, type, active, person_id) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, v := range p.Phones {
		_, err := stmt.ExecContext(ctx, v.Number, v.Type, id)
		if err != nil {
			return fmt.Errorf("can't add phone number: %w", err)
		}
	}
	return nil
}

func(s *Storage) Read(ctx context.Context, data *storage.CrudData) (*entities.Person, error) {
	
	var (
		// FirstName string
		LastName string
		Phone string
		PhoneType entities.PhoneType
		PhoneActive bool
	)
	p := entities.NewPerson()

	var q = `SELECT last_name, number, type, active FROM person LEFT JOIN phone ON 
			phone.person_id = person.id WHERE `
	switch data.SearchComb {
	case storage.Fname:
		p.FirstName = data.FirstName
		q += `first_name = ?`
		rows, err := s.DB.QueryContext(ctx, q, data.FirstName)
		if err != nil {
			return nil, fmt.Errorf("can't query person data: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&LastName, &Phone, &PhoneType, &PhoneActive ); err != nil {
				return nil, fmt.Errorf("can't scan person data to var: %w", err)
			}
			ph := entities.NewPhone(PhoneType, Phone, PhoneActive)
			p.AddPhone(ph)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("got error while getting rows: %w", err)
		}
		p.LastName = LastName
		return p, nil
	case storage.Lname:
		q += `last_name = ?`
	case storage.FullName:
		q += `first_name = ? AND last_name = ?`
	case storage.Phone:
		q += `number = ?`
	case storage.Unknown:
		return nil, storage.ErrNotEnoughSearchCriteria
	}
	return nil, nil
}