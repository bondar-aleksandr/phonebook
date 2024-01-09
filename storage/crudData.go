package storage

import (
	"strings"
)

type CrudComb uint8

const(
	CrudFname CrudComb = iota
	CrudLname
	CrudFullName
	CrudPhone
	CrudAll
	CrudUnknown
)

type CrudData struct {
	FirstName string
	LastName string
	Phone string
	CrudComb CrudComb
}

func NewCrudData(args... string) *CrudData {
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
	searchData := &CrudData{
		FirstName: strings.TrimSpace(FirstName),
		LastName: strings.TrimSpace(LastName),
		Phone: strings.TrimSpace(Phone),
	}
	searchData.setCrudComb()
	return searchData
}

func(s *CrudData) setCrudComb() {
	switch {
	case s.FirstName != "" && s.LastName == "":
		s.CrudComb = CrudFname
	case s.LastName != "" && s.FirstName == "":
		s.CrudComb = CrudLname
	case s.FirstName != "" && s.LastName != "":
		s.CrudComb = CrudFullName
	case s.Phone != "":
		s.CrudComb = CrudPhone
	default:
		s.CrudComb = CrudAll
	}
}