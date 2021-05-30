package app

import (
	Model "RosterService/pkg/Roster/model"
	"errors"
)

var ErrorRosterIdNotExist = errors.New("roster id not found!")

type IRosterDB interface {
	GetAllRosters() ([]Model.Roster, error)
	GetRosterInDBById(id string) (Model.Roster, error)
	DeleteRoster(id string) (string, error)
	GetRosterIdInDBById(id string) (string, error)
}

type IEquipmentService interface {
	GetEquipmentInfo(idEquipment string) (Model.EquipmentDetailedInfo, error)
}

type IUnitService interface {
	GetUnitInfo(idUnit string) (Model.UnitDetailedInfo, error)
}

type RosterApp struct {
	db IRosterDB
	unitService IUnitService
	equipmentService IEquipmentService
}

type EquipmentAppData struct {
	Id string
}

type UnitAppData struct {
	Id string
	Equipments []EquipmentAppData
}

type RosterAppData struct {
	Id string
	Name string
	IdUser string
	Units []UnitAppData
}

func CreateRosterApp(db IRosterDB, unitService IUnitService, equipmentService IEquipmentService) RosterApp {
	return RosterApp{db, unitService, equipmentService}
}

func (app *RosterApp) convertRosterToRosterAppData(roster Model.Roster) RosterAppData {
	rosterApp := RosterAppData {
		roster.Id,
		roster.Name,
		roster.IdUser,
		[]UnitAppData{},
	}

	for i := 0; i < len(roster.Units); i++ {
		unit := roster.Units[i]
		unitApp := UnitAppData {
			unit.Id,
			[]EquipmentAppData{},
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]
			equipmentApp := EquipmentAppData{equipment.Id}
			unitApp.Equipments = append(unitApp.Equipments, equipmentApp)
		}

		rosterApp.Units = append(rosterApp.Units, unitApp)
	}

	return rosterApp
}

func (app *RosterApp) rosterIdExist(id string) bool {
	rosterId, err := app.db.GetRosterIdInDBById(id)
	return (err == nil) && (rosterId != "")
}

func (app *RosterApp) assertRosterNotExist(id string) error {
	if !app.rosterIdExist(id) {
		return ErrorRosterIdNotExist
	}
	return nil
}

func (app *RosterApp) getUnitEqipmentDetailInfo(equipments []Model.Equipment) ([]Model.EquipmentDetailedInfo, error) {
	var equipmentsDeatil []Model.EquipmentDetailedInfo

	for i := 0; i < len(equipments); i++ {
		equipment := equipments[i]
		equipmentDetail, err := app.equipmentService.GetEquipmentInfo(equipment.Id)
		if err != nil {
			return []Model.EquipmentDetailedInfo{}, err
		}
		equipmentsDeatil = append(equipmentsDeatil, equipmentDetail)
	}

	return equipmentsDeatil, nil
}

// TODO БУДУ ИСПОЛЬЗОВАТЬ ДЛЯ ДАЛЬНЕЙШИХ ПРОВЕРОК НА БИЗНЕС ПРАВИЛА ПРИ ДОБАВЛЕНИИ И ИЗМЕНЕНИИ РОСТЕРА
func (app *RosterApp) getFullRosterInformation(roster Model.Roster) (Model.RosterDetailedInfo, error) {
	rosterDetail := Model.RosterDetailedInfo {
		roster.Id,
		roster.Name,
		roster.IdUser,
		[]Model.RosterUnitDetailedInfo{},
	}

	for i := 0; i < len(roster.Units); i++ {
		unit := roster.Units[i]
		unitDetail, err := app.unitService.GetUnitInfo(unit.Id)
		if err != nil {
			return Model.RosterDetailedInfo{}, err
		}

		equipmentsDetail, err := app.getUnitEqipmentDetailInfo(unit.Equipments)
		if err != nil {
			return Model.RosterDetailedInfo{}, err
		}

		rosterUnit := Model.RosterUnitDetailedInfo {unitDetail, equipmentsDetail}
		rosterDetail.Units = append(rosterDetail.Units, rosterUnit)
	}

	return rosterDetail, nil
}

func (app *RosterApp) GetAllRosters() ([]RosterAppData, error) {
	var rosters []RosterAppData
	rostersDB, err := app.db.GetAllRosters()
	if err != nil {
		return rosters, err
	}

	for i := 0; i < len(rostersDB); i++ {
		rosterDB := rostersDB[i]
		rosterApp := app.convertRosterToRosterAppData(rosterDB)
		rosters = append(rosters, rosterApp)
	}

	return rosters, nil
}

func (app *RosterApp) GetRosterById(id string) (RosterAppData, error) {
	roster, err := app.db.GetRosterInDBById(id)
	if err != nil {
		return RosterAppData {}, err
	}
	rosterApp := app.convertRosterToRosterAppData(roster)
	return rosterApp, nil
}

func (app *RosterApp) DeleteById(id string) (string, error) {
	err := app.assertRosterNotExist(id)
	if err != nil {
		return "", err
	}
	deleteId, err := app.db.DeleteRoster(id)
	if err != nil {
		return "", err
	}
	return deleteId, nil
}