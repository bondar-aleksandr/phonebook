package entities

import (
	"errors"
	"fmt"
	"strings"
)

var ErrDublicatePhoneNumber = errors.New("dublicate phone number")

type PhoneType uint8

const (
	Mobile PhoneType = iota
	Work
	Home
	Other
)

func (pt PhoneType) String() string {
    return [...]string{"Mobile", "Work", "Home", "Other"}[pt]
}

type Phone struct {
	Type PhoneType
	Number string
}

func(p *Phone) String() string {
	return fmt.Sprintf("Phone: [Number:%q, Type: %q]", p.Number, p.Type)
}

type Person struct {
	ID int32
	FirstName string
	LastName string
	Notes string
	Phones map[string]*Phone
}

func(p *Person) String() string {
	return fmt.Sprintf("Person: [FirstName: %q, LastName: %q, Notes: %q, Phones: %v]",
	 p.FirstName, p.LastName, p.Notes, p.Phones)
}

func NewPhone(t PhoneType, num string) *Phone {
	return &Phone{
		Type: t,
		Number: num,
	}
}

func NewPerson(args... string) *Person {
	var (
		FirstName string
		LastName string
		Notes string
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
		Notes = args[2]
	default:
	}

	return &Person{
		FirstName: strings.TrimSpace(FirstName),
		LastName: strings.TrimSpace(LastName),
		Notes: strings.TrimSpace(Notes),
		Phones: make(map[string]*Phone),
	}
}

func(p *Person) AddPhone(phone *Phone) error {
	if _, ok := p.Phones[phone.Number]; ok {
		return ErrDublicatePhoneNumber
	}
	p.Phones[phone.Number] = phone
	return nil
}