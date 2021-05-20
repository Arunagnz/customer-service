package postgres

import (
	"fmt"

	"github.com/arunagnz/customer-service/stores"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting database: %w", err)
	}
	return &Store{
		CustomerStore: &stores.CustomerStore{DB: db},
	}, nil
}

type Store struct {
	*stores.CustomerStore
}
