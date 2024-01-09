package storage

import (
	"context"
	"errors"
	"github.com/bondar-aleksandr/phonebook/entities"
)

var ErrNoPersonExists = errors.New("no persons exists")

type Storage interface {
	Create(context.Context, *entities.Person) error
	Read(context.Context, *CrudData) ([]*entities.Person, error)
	// Update(context.Context, *entities.Person) error
	Delete(context.Context, *CrudData) (int64, error)
	Reset(context.Context) error
	Populate(context.Context) error
}