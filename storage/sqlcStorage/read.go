package sqlcstorage

import (
	"context"
	"fmt"
	"github.com/bondar-aleksandr/phonebook/db_access"
	"github.com/bondar-aleksandr/phonebook/entities"
	"github.com/bondar-aleksandr/phonebook/storage"
	_ "github.com/go-sql-driver/mysql"
)


// Read CRUD operation.
// TODO: refactor if possible, don't to repeat the same actions for every case!
func(s *Storage) Read(ctx context.Context, cd *storage.CrudData) ([]*entities.Person, error) {
	var persons []*entities.Person
	switch cd.CrudComb {
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