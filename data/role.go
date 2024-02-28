package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Role struct {
	ID          int           `db:"id"`
	Name        string        `db:"name"`
	Alignment   string        `db:"alignment"`
	AbilityIDs  pq.Int32Array `db:"ability_ids"`
	PassiveIDs  pq.Int32Array `db:"passive_ids"`
}

type ComplexRole struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Alignment   string     `db:"alignment"`
	Abilities   []*Ability `db:"abilities"`
	Passives    []*Passive `db:"passives"`
}

type RoleModel struct {
	*sqlx.DB
}

func (rm *RoleModel) Get(id int) (*Role, error) {
	query := `SELECT * FROM roles WHERE id = $1`
	var role Role
	err := rm.DB.Get(&role, query, id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (rm *RoleModel) GetByName(name string) (*Role, error) {
	var role Role
	err := rm.DB.Get(&role, "SELECT * FROM roles WHERE name ILIKE $1", name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (rm *RoleModel) GetComplex(id int) (*ComplexRole, error) {
	query := `SELECT r.id, r.name, r.alignment, 
  (SELECT array_agg(a) FROM abilities a WHERE a.id = ANY(r.ability_ids)) as abilities,
  (SELECT array_agg(p) FROM passives p WHERE p.id = ANY(r.passive_ids)) as passives
  FROM roles r WHERE r.id = $1`
	var complexRole ComplexRole
	err := rm.DB.Get(&complexRole, query, id)
	if err != nil {
		return nil, err
	}
	return &complexRole, nil
}

func (rm *RoleModel) GetComplexByName(name string) (*ComplexRole, error) {
	query := `SELECT r.id, r.name, r.alignment, 
  (SELECT array_agg(a) FROM abilities a WHERE a.id = ANY(r.ability_ids)) as abilities,
  (SELECT array_agg(p) FROM passives p WHERE p.id = ANY(r.passive_ids)) as passives
  FROM roles r WHERE r.name ILIKE $1`
	var complexRole ComplexRole
	err := rm.DB.Get(&complexRole, query, name)
	if err != nil {
		return nil, err
	}
	return &complexRole, nil
}

func (rm *RoleModel) GetAll() ([]*Role, error) {
	var roles []*Role
	err := rm.DB.Select(&roles, "SELECT * FROM roles")
	if err != nil {
		return nil, err
	}
	return roles, nil
}
