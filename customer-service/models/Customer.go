package models

import "github.com/google/uuid"

type Customer struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	PhoneNumber string    `db:"phonenumber"`
	Address     string    `db:"address"`
}

type CustomerStore interface {
	Customer(id uuid.UUID) (Customer, error)
	Customers() ([]Customer, error)
	CreateCustomer(u *Customer) error
	UpdateCustomer(u *Customer) error
	DeleteCustomer(id uuid.UUID) error
}

type Store interface {
	CustomerStore
}
