package main

import "testing"

func TestReadRoleCSV(t *testing.T) {
	filepath := "./Betrayal Tabletop Online (BTO) Information Spreadsheet [2024] - ROLES.csv"

	roles, err := ReadRoleCSV(filepath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(roles[0])
}
