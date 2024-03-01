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

func ReadItemCSV(filepath string) ([]data.Item, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	items := []data.Item{}
	i := 0
	for {
		var item data.Item
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// skip first line of csv
		if i == 0 || i == 1 {
			i++
			continue
		}
		item.Rarity = strings.TrimSpace(record[1])
		item.Name = strings.TrimSpace(record[2])
		cost := -1
		if record[3] != "X" {
			cost, err = strconv.Atoi(record[3])
			if err != nil {
				return nil, err
			}
		}
		categories := strings.Split(record[4], "/")
		for _, c := range categories {
			item.Categories = append(item.Categories, strings.TrimSpace(c))
		}
		item.Description = strings.TrimSpace(record[5])
		item.Cost = cost
		items = append(items, item)
		i++
	}

	return items, nil
}

func ReadStatusCSV(filepath string) ([]data.Status, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	statuses := []data.Status{}
	i := 0
	for {
		var status data.Status
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// skip first line of csv
		if i == 0 || i == 1 {
			i++
			continue
		}
		status.Name = strings.TrimSpace(record[1])
		status.Description = strings.TrimSpace(record[2])
		statuses = append(statuses, status)
		i++
	}

	return statuses, nil
}
