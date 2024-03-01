package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func TestReadRoleCSV(t *testing.T) {
	dbStr := os.Getenv("DATABASE_URL")

	db, err := sqlx.Open("postgres", dbStr)
	if err != nil {
		t.Fatal(err)
	}

	filepath := "./Betrayal Tabletop Online (BTO) Information Spreadsheet [2024] - ROLES.csv"
	entries, err := ReadRoleCSV(filepath)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range entries {
		fmt.Println("Inserting", v.R.Name)
		role := v.R
		abilities := v.A
		passives := v.P

		// Insert abilities into psql database
		abResults := pq.Int32Array{}
		for _, a := range abilities {
			// insert the ability in psql
			_, err := db.NamedExec("INSERT INTO abilities (name, charges, description, role_specific, any_ability, rarity, categories) VALUES (:name, :charges, :description, :role_specific, :any_ability, :rarity, :categories)", a)
			if err != nil {
				t.Fatal(err)
			}
			// get all the ids of the abilities
			var id int32
			err = db.Get(&id, "SELECT id FROM abilities WHERE name = $1", a.Name)
			if err != nil {
				t.Fatal(err)
			}
			abResults = append(abResults, id)
		}

		// Insert passives into psql database
		pResults := pq.Int32Array{}
		for _, p := range passives {
			_, err := db.NamedExec("INSERT INTO passives (name, description) VALUES (:name, :description)", p)
			if err != nil {
				t.Fatal(err)
			}
			var id int32
			err = db.Get(&id, "SELECT id FROM passives WHERE name = $1", p.Name)
			if err != nil {
				t.Fatal(err)
			}
			pResults = append(pResults, id)
		}
		role.AbilityIDs = append(role.AbilityIDs, abResults...)
		role.PassiveIDs = append(role.PassiveIDs, pResults...)

		// Insert role into psql database
		_, err := db.NamedExec("INSERT INTO roles (name, alignment, ability_ids, passive_ids) VALUES (:name, :alignment, :ability_ids, :passive_ids)", role)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestReadItemCSV(t *testing.T) {
	dbStr := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dbStr)
	if err != nil {
		t.Fatal(err)
	}
	filepath := "./Betrayal Tabletop Online (BTO) Information Spreadsheet [2024] - Items.csv"
	entries, err := ReadItemCSV(filepath)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range entries {
		fmt.Println("Inserting", v.Name)
		item := v
		// Insert item into psql database
		_, err := db.NamedExec("INSERT INTO items (name, description, rarity, cost, categories) VALUES (:name, :description, :rarity, :cost, :categories)", item)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestReadStatusCSV(t *testing.T) {
	dbStr := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dbStr)
	if err != nil {
		t.Fatal(err)
	}
	filepath := "./Betrayal Tabletop Online (BTO) Information Spreadsheet [2024] - Statuses.csv"
	entries, err := ReadStatusCSV(filepath)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range entries {
		fmt.Println("Inserting", v.Name)
		status := v
		// Insert status into psql database
		_, err := db.NamedExec("INSERT INTO statuses (name, description) VALUES (:name, :description)", status)
		if err != nil {
			t.Fatal(err)
		}
	}
}
