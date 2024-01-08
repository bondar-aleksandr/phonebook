package phonebook

import (
	"context"
	"github.com/bondar-aleksandr/phonebook/entities"
	"github.com/bondar-aleksandr/phonebook/storage"
)

type PhoneBook struct {
	storage storage.Storage
}

func New(s storage.Storage) *PhoneBook {
	return &PhoneBook{
		storage: s,
	}
}

func (p *PhoneBook) Add(ctx context.Context, person *entities.Person) error {
	return p.storage.Create(ctx, person)
}

func (p *PhoneBook) Search(ctx context.Context, s *entities.SearchData) ([]*entities.Person, error) {
	return p.storage.Read(ctx, s)
}