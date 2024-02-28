package data

import "github.com/jmoiron/sqlx"

type Player struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	GameID       string `db:"game_id"`
	RoleID       int    `db:"role_id"`
	Alive        bool   `db:"alive"`
	Seat         int    `db:"seat"`
	Luck         int    `db:"luck"`
	LuckModifier int    `db:"luck_modifier"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

// ComplexPlayer is a player with a role
type ComplexPlayer struct {
	P Player
	R Role
}

type PlayerModel struct {
	DB *sqlx.DB
}

func (m *PlayerModel) GetByGameID(gameID string) ([]*Player, error) {
	players := []*Player{}
	err := m.DB.Select(&players, "SELECT * FROM players WHERE game_id = $1", gameID)
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (m *PlayerModel) GetByID(id int) (*Player, error) {
	player := &Player{}
	err := m.DB.Get(player, "SELECT * FROM players WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (m *PlayerModel) GetByName(name string) (*Player, error) {
	player := &Player{}
	err := m.DB.Get(player, "SELECT * FROM players WHERE name ILIKE $1", name)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (m *PlayerModel) GetByGameIDAndName(gameID string, name string) (*Player, error) {
	player := &Player{}
	err := m.DB.Get(player, "SELECT * FROM players WHERE game_id = $1 AND name ILIKE $2", gameID, name)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (m *PlayerModel) Create(player *Player) error {
	_, err := m.DB.NamedExec("INSERT INTO players (name, game_id, role_id, alive, seat, luck, luck_modifier) VALUES (:name, :game_id, :role_id, :alive, :seat, :luck, :luck_modifier)", player)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) Update(player *Player) error {
	_, err := m.DB.NamedExec("UPDATE players SET name = :name, game_id = :game_id, role_id = :role_id, alive = :alive, seat = :seat, luck = :luck, luck_modifier = :luck_modifier WHERE id = :id", player)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM players WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PlayerModel) GetRole(roleID int) (*Role, error) {
	role := &Role{}
	err := m.DB.Get(role, "SELECT * FROM roles WHERE id = $1", roleID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (m *PlayerModel) GetComplexByGameID(gameID string) ([]*ComplexPlayer, error) {
	players := []*ComplexPlayer{}
	simplePlayers, err := m.GetByGameID(gameID)
	if err != nil {
		return nil, err
	}

	for _, player := range simplePlayers {
		role, err := m.GetRole(player.RoleID)
		if err != nil {
			return nil, err
		}
		players = append(players, &ComplexPlayer{P: *player, R: *role})
	}

	return players, nil
}
