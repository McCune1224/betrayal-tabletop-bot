package data

import "github.com/jmoiron/sqlx"

type Status struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type StatusModel struct {
	*sqlx.DB
}

func (m *StatusModel) GetAll() ([]Status, error) {
	var statuses []Status
	query := `SELECT * FROM statuses`
	err := m.Select(&statuses, query)
	return statuses, err
}

func (m *StatusModel) GetByID(id int) (Status, error) {
	var status Status
	query := `SELECT * FROM statuses WHERE id = $1`
	err := m.Get(&status, query, id)
	return status, err
}

func (m *StatusModel) GetByName(name string) (Status, error) {
	var status Status
	query := `SELECT * FROM statuses WHERE name = $1`
	err := m.Get(&status, query, name)
	return status, err
}
