package storage

import (
	"context"
	"errors"
	"strings"
	"github.com/bondar-aleksandr/phonebook/entities"
)

var ErrNoPersonExists = errors.New("no persons exists")
var ErrNotEnoughSearchCriteria = errors.New("not enough search criteria specified")

type SearchComb uint8

const(
	Fname SearchComb = iota
	Lname
	FullName
	Phone
	Unknown
)

type SearchData struct {
	FirstName string
	LastName string
	Phone string
	SearchComb SearchComb
}

func NewSearchData(args... string) *SearchData {
	var (
		FirstName string
		LastName string
		Phone string
	)
	switch len(args) {
	case 1:
		FirstName = args[0]
	case 2:
		FirstName = args[0]
		LastName = args[1]
	case 3:
		FirstName = args[0]
		LastName = args[1]
		Phone = args[2]
	}
	searchData := &SearchData{
		FirstName: strings.TrimSpace(FirstName),
		LastName: strings.TrimSpace(LastName),
		Phone: strings.TrimSpace(Phone),
	}
	searchData.setSearchType()
	return searchData
}

func(s *SearchData) setSearchType() {
	switch {
	case s.FirstName != "" && s.LastName == "":
		s.SearchComb = Fname
	case s.LastName != "" && s.FirstName == "":
		s.SearchComb = Lname
	case s.FirstName != "" && s.LastName != "":
		s.SearchComb = FullName
	case s.Phone != "":
		s.SearchComb = Phone
	default:
		s.SearchComb = Unknown
	}
}

type Storage interface {
	Create(context.Context, *entities.Person) error
	Read(context.Context, *SearchData) (*entities.Person, error)
	// Update(context.Context, *entities.Person) error
	// Delete(context.Context, *entities.Person) error
	// ReadAll(context.Context) error
}