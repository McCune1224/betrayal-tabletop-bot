package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mccune1224/betrayal-tabletop-bot/data"
)

func ReadRoleCSV(filepath string) ([]SheetRole, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	chunks := [][][]string{}
	currChunk := [][]string{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			chunks = append(chunks, currChunk)
			break
		}
		if err != nil {
			return nil, err
		}
		if record[1] == "" {
			chunks = append(chunks, currChunk)
			currChunk = [][]string{}
		} else {
			currChunk = append(currChunk, record)
		}
	}
	var roles []SheetRole

	for _, chunk := range chunks[1:] {
		role := ParseChunk(chunk)
		roles = append(roles, role)
	}

	return roles, nil
}

func ParseChunk(chunk [][]string) SheetRole {
	var role SheetRole
	role.R.Name = strings.TrimSpace(chunk[1][1])
	role.R.Alignment = strings.TrimSpace(chunk[3][1])
	abStartIdx := 5
	for strings.TrimSpace(chunk[abStartIdx][1]) != "Passives:" {
		role.A = append(role.A, ParseAbility(chunk[abStartIdx], role.R.Name))
		abStartIdx++
	}

	for _, p := range chunk[abStartIdx+1:] {
		role.P = append(role.P, data.Passive{Name: strings.TrimSpace(p[1]), Description: strings.TrimSpace(p[2])})
	}
	return role
}

func ParseAbility(line []string, roleName string) data.Ability {
	a := data.Ability{}
	a.Name = strings.TrimSpace(line[1])
	charge, _ := strconv.Atoi(line[2])
	a.Charges = charge

	switch line[3] {
	case "*":
		a.AnyAbility = true
		a.RoleSpecific = ""
		a.Rarity = line[6]
	case "^":
		a.AnyAbility = true
		a.RoleSpecific = roleName
		a.Rarity = line[6]
	case "":
		a.AnyAbility = false
		a.RoleSpecific = roleName
		a.Rarity = ""
	}
	a.Description = strings.TrimSpace(line[4])
	categories := strings.Split(line[5], "/")
	for _, c := range categories {
		a.Categories = append(a.Categories, strings.TrimSpace(c))
	}

	return a
}

type SheetRole struct {
	R data.Role
	A []data.Ability
	P []data.Passive
}
