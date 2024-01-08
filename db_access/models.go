// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db_access

import (
	"database/sql"
)

type Person struct {
	ID        int32
	FirstName string
	LastName  string
	Notes     sql.NullString
	Created   sql.NullTime
	Modified  sql.NullTime
}

type Phone struct {
	ID          int32
	PhoneNumber string
	PhoneType   int32
	PersonID    int32
}
