package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Ability struct {
	ID           int            `db:"id"`
	Name         string         `db:"name"`
	Description  string         `db:"description"`
	Charges      int            `db:"charges"`
	AnyAbility   bool           `db:"any_ability"`
	RoleSpecific string         `db:"role_specific"`
	Rarity       string         `db:"rarity"`
	Categories   pq.StringArray `db:"categories"`
}

// psql statement to add categories to ability
// ALTER TABLE abilities ADD COLUMN categories text[] DEFAULT '{}';

type AbilityModel struct {
	*sqlx.DB
}

func (am *AbilityModel) Get(id int) (*Ability, error) {
	query := `SELECT * FROM abilities WHERE id = $1`
	var ability Ability
	err := am.DB.Get(&ability, query, id)
	if err != nil {
		return nil, err
	}
	return &ability, nil
}

func (am *AbilityModel) GetByName(name string) (*Ability, error) {
	query := `SELECT * FROM abilities WHERE name ILIKE $1`
	var ability Ability
	err := am.DB.Get(&ability, query, name)
	if err != nil {
		return nil, err
	}
	return &ability, nil
}

func (am *AbilityModel) GetAll() ([]*Ability, error) {
	query := `SELECT * FROM abilities`
	var abilities []*Ability
	err := am.Select(&abilities, query)
	if err != nil {
		return nil, err
	}
	return abilities, nil
}

func (am *AbilityModel) GetRandomByRarity(rarity string) (*Ability, error) {
	query := `SELECT * FROM abilities WHERE rarity ILIKE $1 ORDER BY random() LIMIT 1`
	var ability Ability
	err := am.DB.Get(&ability, query, rarity)
	if err != nil {
		return nil, err
	}
	return &ability, nil
}
