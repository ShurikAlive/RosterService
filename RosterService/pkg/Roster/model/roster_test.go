package model

import "testing"

func createValidRosterInputData() RosterInput {
	roster := RosterInput{
		"test roster id",
		"test roster",
		"test user id",
		[]UnitInput{{"test unit id 1", []EquipmentInput{} },
			{"test unit id 2", []EquipmentInput{} },
			{"test unit id 3", []EquipmentInput{} },
			{"test unit id 4", []EquipmentInput{} },
		},
	}

	return roster
}

func TestCreateRoster(t *testing.T) {
	rosterIn := createValidRosterInputData()
	_, err := CreateRoster(rosterIn)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidNameRoster(t *testing.T) {
	rosterIn := createValidRosterInputData()
	rosterIn.Name = ""
	_, err := CreateRoster(rosterIn)
	if err != ErrorRosterName {
		t.Fatal(err)
	}
}

func TestInvalidUserIdRoster(t *testing.T) {
	rosterIn := createValidRosterInputData()
	rosterIn.IdUser = ""
	_, err := CreateRoster(rosterIn)
	if err != ErrorRosterUserId {
		t.Fatal(err)
	}
}

func TestInvalidCountUnitInRoster(t *testing.T) {
	rosterIn := createValidRosterInputData()
	i := 1
	rosterIn.Units = rosterIn.Units[:i+copy(rosterIn.Units[i:], rosterIn.Units[i+1:])]// Удаляем i элемент
	_, err := CreateRoster(rosterIn)
	if err != ErrorCountUnitInRoster {
		t.Fatal(err)
	}
}

func TestInvalidMaxCountUnitInRoster(t *testing.T) {
	rosterIn := createValidRosterInputData()
	unit := UnitInput{"test unit id 1", []EquipmentInput{}, }
	rosterIn.Units = append(rosterIn.Units, unit)
	_, err := CreateRoster(rosterIn)
	if err != ErrorCountUnitInRoster {
		t.Fatal(err)
	}
}

func TestInvalidDoubleUnitInRoster(t *testing.T) {
	rosterIn := createValidRosterInputData()
	rosterIn.Units[1].Id = rosterIn.Units[2].Id
	_, err := CreateRoster(rosterIn)
	if err != ErrorDoubleUnit {
		t.Fatal(err)
	}
}