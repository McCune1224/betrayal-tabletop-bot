package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

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
	role.R.Name = chunk[1][1]
	role.R.Alignment = chunk[3][1]

	log.Println(role.R.Name, role.R.Alignment)
	return role
}

type SheetRole struct {
	R data.Role
	A []data.Ability
	P []data.Passive
}
