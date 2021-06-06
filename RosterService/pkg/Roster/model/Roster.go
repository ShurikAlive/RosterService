package model

import "errors"

var ErrorCountUnitInRoster = errors.New("error count unit in roster!")
var ErrorRosterName = errors.New("error roster name!")
var ErrorDoubleUnit = errors.New("error double unit!")
var ErrorRosterUserId = errors.New("error roster user id!")

var COUNT_UNIT_IN_ROSTER = 4

type Equipment struct {
	Id string
}

type Unit struct {
	Id string
	Equipments []Equipment
}

type Roster struct {
	Id string
	Name string
	IdUser string
	Units []Unit
}

type EquipmentInput struct {
	Id string
}

type UnitInput struct {
	Id string
	Equipments []EquipmentInput
}

type RosterInput struct {
	Id     string
	Name   string
	IdUser string
	Units  []UnitInput
}

func createRoster(inRoster RosterInput) Roster {
	roster := Roster {
		inRoster.Id,
		inRoster.Name,
		inRoster.IdUser,
		[]Unit{},
	}

	for i := 0; i < len(inRoster.Units); i++ {
		unit := inRoster.Units[i]
		unitApp := Unit {
			unit.Id,
			[]Equipment{},
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]
			equipmentApp := Equipment{equipment.Id}
			unitApp.Equipments = append(unitApp.Equipments, equipmentApp)
		}

		roster.Units = append(roster.Units, unitApp)
	}

	return roster
}

func countUnitInRoster(roster RosterInput) error {
	if len(roster.Units) != COUNT_UNIT_IN_ROSTER {
		return ErrorCountUnitInRoster
	}
	return nil
}

func nameRosterValid(roster RosterInput) error {
	if roster.Name == "" {
		return ErrorRosterName
	}
	return nil
}

func userIdRosterValid(roster RosterInput) error {
	if roster.IdUser == "" {
		return ErrorRosterUserId
	}
	return nil
}

func doubleUnitInRoster(roster RosterInput) error {
	unitCount := make(map[string]int32)
	for i := 0; i < len(roster.Units); i++ {
		_, unitExistInMap := unitCount[roster.Units[i].Id]
		if unitExistInMap {
			return ErrorDoubleUnit
		} else {
			unitCount[roster.Units[i].Id] = 1
		}
	}
	return nil
}

func CreateRoster(inRoster RosterInput) (Roster, error) {
	err := countUnitInRoster(inRoster)
	if err != nil {
		return Roster{}, err
	}
	err = nameRosterValid(inRoster)
	if err != nil {
		return Roster{}, err
	}
	err = doubleUnitInRoster(inRoster)
	if err != nil {
		return Roster{}, err
	}
	err = userIdRosterValid(inRoster)
	if err != nil {
		return Roster{}, err
	}
	return createRoster(inRoster), nil
}
