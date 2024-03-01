package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Item struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	Rarity      string         `db:"rarity"`
	Cost        int            `db:"cost"`
	Categories  pq.StringArray `db:"categories"`
}

type ItemModel struct {
	*sqlx.DB
}

func (m *ItemModel) Get(id int) (*Item, error) {
	var item Item
	err := m.DB.Get(&item, "SELECT * FROM items WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (m *ItemModel) GetAll() ([]*Item, error) {
	var items []*Item
	err := m.DB.Select(&items, "SELECT * FROM items")
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (m *ItemModel) GetByName(name string) (*Item, error) {
	var item Item
	err := m.DB.Get(&item, "SELECT * FROM items WHERE name ILIKE $1", name)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
