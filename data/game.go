package data

import "github.com/jmoiron/sqlx"

type Game struct {
	ID          int    `json:"id" db:"id"`
	GameID      string `json:"game_id" db:"game_id"`
	PlayerCount int    `json:"player_count" db:"player_count"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type GameModel struct {
	*sqlx.DB
}

func (gm *GameModel) GetByID(id int) (*Game, error) {
	var game Game
	err := gm.Get(&game, "SELECT * FROM games WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (gm *GameModel) GetAll() ([]Game, error) {
	var games []Game
	err := gm.Select(&games, "SELECT * FROM games")
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (gm *GameModel) GetByGameID(gID string) (*Game, error) {
	var game Game
	err := gm.Get(&game, "SELECT * FROM games WHERE game_id = $1", gID)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (gm *GameModel) InsertGame(gameID string, playerCount int) (*Game, error) {
	_, err := gm.Exec("INSERT INTO games (game_id, player_count) VALUES ($1, $2)", gameID, playerCount)
	if err != nil {
		return nil, err
	}
	return gm.GetByGameID(gameID)
}

func (gm *GameModel) Update(game *Game) error {
	_, err := gm.Exec("UPDATE games SET game_id = $1, player_count = $2 WHERE id = $3", game.GameID, game.PlayerCount, game.ID)
	if err != nil {
		return err
	}
	return nil
}

func (gm *GameModel) UpdatePlayerCount(gID string, playerCount int) error {
	_, err := gm.Exec("UPDATE games SET player_count = $1 WHERE game_id = $2", playerCount, gID)
	if err != nil {
		return err
	}
	return nil
}

func (gm *GameModel) DeleteGame(gID string) error {
	_, err := gm.Exec("DELETE FROM games WHERE game_id = $1", gID)
	if err != nil {
		return err
	}

	// delete all players in the game with the given game_id
	_, err = gm.Exec("DELETE FROM players WHERE game_id = $1", gID)
	if err != nil {
		return err
	}
	return nil
}
