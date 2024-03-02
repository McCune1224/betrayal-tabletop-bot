package data

import (
	"github.com/jmoiron/sqlx"
)

type Passive struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type PassiveModel struct {
	*sqlx.DB
}

func (pm *PassiveModel) Get(id int) (*Passive, error) {
	query := `SELECT * FROM passives WHERE id = $1`
	var passive Passive
	err := pm.DB.Get(&passive, query, id)
	if err != nil {
		return nil, err
	}
	return &passive, nil
}

func (pm *PassiveModel) GetByName(name string) (*Passive, error) {
	query := `SELECT * FROM passives WHERE name = $1`
	var passive Passive
	err := pm.DB.Get(&passive, query, name)
	if err != nil {
		return nil, err
	}
	return &passive, nil
}

func (pm *PassiveModel) GetAll() ([]*Passive, error) {
	query := `SELECT * FROM passives`
	var passives []*Passive
	err := pm.Select(&passives, query)
	if err != nil {
		return nil, err
	}
	return passives, nil
}
