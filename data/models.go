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
}

func NewModels(db *sqlx.DB) *Models {
	return &Models{
		Games:     GameModel{DB: db},
		Players:   PlayerModel{DB: db},
		Roles:     RoleModel{DB: db},
		Abilities: AbilityModel{DB: db},
		Passives:  PassiveModel{DB: db},
	}
}
