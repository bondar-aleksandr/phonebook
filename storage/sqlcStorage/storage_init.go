package sqlcstorage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
	"github.com/bondar-aleksandr/phonebook/db_access"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	DB *sql.DB
	queries *db_access.Queries
}

// New creates new instance of storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("mysql", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	s := &Storage{DB: db}
	if err = s.init(context.TODO()); err != nil {
		return nil, fmt.Errorf("can't initialize database: %w", err)
	}
	s.queries = db_access.New(db)
	return s, nil
}

// Init executes commands from sql_init.sql file. Creates database
// and tables.
func(s *Storage) init(ctx context.Context) error {
	sqlInit, err := os.ReadFile("sql/init.sql")
	if err != nil {
		return  fmt.Errorf("can't open sql init file: %w", err)
	}
	_, err = s.DB.ExecContext(ctx, string(sqlInit))
	if err != nil {
		return  fmt.Errorf("can't execute sql init file: %w", err)
	}
	return nil
}

// Reset executes commands from sql_reset.sql file. Drops database.
func(s *Storage) Reset(ctx context.Context) error {
	sqlreset, err := os.ReadFile("sql/reset.sql")
	if err != nil {
		return  fmt.Errorf("can't open sql reset file: %w", err)
	}
	_, err = s.DB.ExecContext(ctx, string(sqlreset))
	if err != nil {
		return  fmt.Errorf("can't execute sql reset file: %w", err)
	}
	return nil
}

func(s *Storage) Populate(ctx context.Context) error {
	sqlreset, err := os.ReadFile("sql/populate.sql")
	if err != nil {
		return  fmt.Errorf("can't open sql populate file: %w", err)
	}
	_, err = s.DB.ExecContext(ctx, string(sqlreset))
	if err != nil {
		return  fmt.Errorf("can't execute sql populate file: %w", err)
	}
	return nil
}
