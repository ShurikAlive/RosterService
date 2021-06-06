package transport

import (
	App "RosterService/pkg/roster/app"
	"encoding/json"
)

type JsonFormatter struct {

}

type Equipment struct {
	// ID equipment
	Id string `json:"id"`
}

type Unit struct {
	// ID unit
	Id string `json:"id"`

	Equipments []Equipment `json:"equipments"`
}

type Roster struct {
	// ID roster
	Id string `json:"id"`
	// NAME roster
	Name string `json:"name"`
	// ID user
	IdUser string `json:"idUser"`

	Units []Unit `json:"units"`
}

type RosterEdit struct {
	// NAME roster
	Name string `json:"name"`
	// ID user
	IdUser string `json:"idUser"`

	Units []Unit `json:"units"`
}

func CreateJSONFormatter() JsonFormatter {
	return JsonFormatter{}
}

func (formatter *JsonFormatter) ConvertIdToJSON (id string) string {
	return "\"" + id + "\""
}

func (formatter *JsonFormatter) createEditRosterAppData(roster RosterEdit) App.EditRosterAppData {
	rosterApp := App.EditRosterAppData{
		roster.Name,
		roster.IdUser,
		[]App.UnitAppData{},
	}

	for i := 0; i < len(roster.Units); i++ {
		unit := roster.Units[i]
		unitApp := App.UnitAppData {
			unit.Id,
			[]App.EquipmentAppData{},
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]
			equipmentApp := App.EquipmentAppData{equipment.Id}
			unitApp.Equipments = append(unitApp.Equipments, equipmentApp)
		}

		rosterApp.Units = append(rosterApp.Units, unitApp)
	}

	return rosterApp
}

func (formatter *JsonFormatter) ConvertJsonToRocterEditAppData (rosterJson []byte) (App.EditRosterAppData, error) {
	var msg RosterEdit
	err := json.Unmarshal(rosterJson, &msg)
	if err != nil {
		return App.EditRosterAppData{}, err
	}

	return formatter.createEditRosterAppData(msg), nil
}

func (formatter *JsonFormatter) createRoster(roster App.RosterAppData) Roster {
	rosterJson := Roster{
		roster.Id,
		roster.Name,
		roster.IdUser,
		[]Unit{},
	}

	for i := 0; i < len(roster.Units); i++ {
		unit := roster.Units[i]
		unitJson := Unit {
			unit.Id,
			[]Equipment{},
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]
			equipmentJson := Equipment{equipment.Id}
			unitJson.Equipments = append(unitJson.Equipments, equipmentJson)
		}

		rosterJson.Units = append(rosterJson.Units, unitJson)
	}

	return rosterJson
}

func (formatter *JsonFormatter) ConvertRosterAppDataToJSON(roster App.RosterAppData)  ([]byte, error) {
	rosterJson := formatter.createRoster(roster)
	b, err := json.Marshal(rosterJson)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (formatter *JsonFormatter) createAllRosters(rosters []App.RosterAppData) []Roster {
	var rostersJson []Roster
	for i := 0; i < len(rosters); i++ {
		roster := formatter.createRoster(rosters[i])
		rostersJson = append(rostersJson, roster)
	}
	return rostersJson
}

func (formatter *JsonFormatter) ConvertAllRostersAppDataToJSON(rosters []App.RosterAppData)  ([]byte, error) {
	rostersJson := formatter.createAllRosters(rosters)
	b, err := json.Marshal(rostersJson)
	if err != nil {
		return nil, err
	}

	return b, nil
}
