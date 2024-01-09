package sqlcstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bondar-aleksandr/phonebook/db_access"
	"github.com/bondar-aleksandr/phonebook/entities"
	_ "github.com/go-sql-driver/mysql"
)


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
