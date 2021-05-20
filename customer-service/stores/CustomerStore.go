package stores

import (
	"fmt"
	"github.com/arunagnz/customer-service/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewCustomerStore(db *sqlx.DB) *CustomerStore {
	return &CustomerStore{
		DB: db,
	}
}

type CustomerStore struct {
	*sqlx.DB
}

func (s *CustomerStore) Customer(id uuid.UUID) (models.Customer, error) {
	var c models.Customer
	if err := s.Get(&c, `SELECT * from customers WHERE id = $1`, id); err != nil {
		return models.Customer{}, fmt.Errorf("error getting Customer: %w", err)
	}
	return c, nil
}

func (s *CustomerStore) Customers() ([]models.Customer, error) {
	var cc []models.Customer
	if err := s.Select(&cc, `SELECT * from customers`); err != nil {
		return []models.Customer{}, fmt.Errorf("error getting Customers: %w", err)
	}
	return cc, nil
}

func (s *CustomerStore) CreateCustomer(c *models.Customer) error {
	if err := s.Get(c, `INSERT INTO customers VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		c.ID,
		c.Name,
		c.Email,
		c.Password,
		c.PhoneNumber,
		c.Address); err != nil {
		return fmt.Errorf("error creating Customer: %w", err)
	}
	return nil
}

func (s *CustomerStore) UpdateCustomer(c *models.Customer) error {
	if err := s.Get(c, `UPDATE customers SET name = $1, email = $2, password = $3, phonenumber = $4, address = $5 WHERE id = $6 RETURNING *`,
		c.Name,
		c.Email,
		c.Password,
		c.PhoneNumber,
		c.Address,
		c.ID); err != nil {
		return fmt.Errorf("error updating Customer: %w", err)
	}
	return nil
}

func (s *CustomerStore) DeleteCustomer(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM Customers WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting Customer: %w", err)
	}
	return nil
}
