package storage

import (
	"context"
	"errors"
	"github.com/bondar-aleksandr/phonebook/entities"
)

var ErrNoPersonExists = errors.New("no persons exists")
var ErrNotEnoughSearchCriteria = errors.New("not enough search criteria specified")



type Storage interface {
	Create(context.Context, *entities.Person) error
	Read(context.Context, *entities.SearchData) ([]*entities.Person, error)
	// Update(context.Context, *entities.Person) error
	// Delete(context.Context, *entities.Person) error
	// ReadAll(context.Context) error
}