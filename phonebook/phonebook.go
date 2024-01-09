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

func (p *PhoneBook) Search(ctx context.Context, s *storage.CrudData) ([]*entities.Person, error) {
	return p.storage.Read(ctx, s)
}

func (p *PhoneBook) Delete(ctx context.Context, s *storage.CrudData) (int64, error) {
	return p.storage.Delete(ctx, s)
}

func(p *PhoneBook) Reset(ctx context.Context) error {
	return p.storage.Reset(ctx)
}

func(p *PhoneBook) Populate(ctx context.Context) error {
	return p.storage.Populate(ctx)
}