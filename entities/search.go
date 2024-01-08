package entities

import (
	"strings"
)

type Specification interface {
	IsSatisfied() bool
}

type SearchComb uint8

const(
	SearchFname SearchComb = iota
	SearchLname
	SearchFullName
	SearchPhone
	SearchAll
	SearchUnknown
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
		s.SearchComb = SearchFname
	case s.LastName != "" && s.FirstName == "":
		s.SearchComb = SearchLname
	case s.FirstName != "" && s.LastName != "":
		s.SearchComb = SearchFullName
	case s.Phone != "":
		s.SearchComb = SearchPhone
	default:
		s.SearchComb = SearchAll
	}
}