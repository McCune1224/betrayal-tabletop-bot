package data

import (
	"github.com/jmoiron/sqlx"
)

type Models struct {
	Games     GameModel
	Players   PlayerModel
	Roles     RoleModel
	Abilities AbilityModel
	Passives  PassiveModel
	Items     ItemModel
	Statuses  StatusModel
}

func NewModels(db *sqlx.DB) *Models {
	return &Models{
		Games:     GameModel{DB: db},
		Players:   PlayerModel{DB: db},
		Roles:     RoleModel{DB: db},
		Abilities: AbilityModel{DB: db},
		Passives:  PassiveModel{DB: db},
		Items:     ItemModel{DB: db},
		Statuses:  StatusModel{DB: db},
	}
}
