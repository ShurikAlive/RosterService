package app

import (
	Model "RosterService/pkg/roster/model"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
)

var ErrorRosterIdNotExist = errors.New("roster id not found!")
var ErrorRosterExist = errors.New("roster exist!")
var ErrorAccessUserRoster = errors.New("error access user roster!")

type RosterRepository interface {
	Lock() error
	Unlock() error

	GetAllRosters() ([]Model.Roster, error)
	GetRosterById(id string) (Model.Roster, error)
	DeleteRoster(id string) (string, error)
	GetRosterIdById(id string) (string, error)
	GetRosterIdByParams(params RequiredParameters) (string, error)
	InsertNewRoster(roster Model.Roster) (string, error)
	UpdateRoster(roster Model.Roster) (string, error)
	GetOwnerRoster(id string) (string, error)
	GetRostersIdsByUnitId(unitId string) ([]string, error)
	UpdateRosterStatus(id string, status int32) error
	DeleteUnitInRoster(idRoster string, idUnit string) error
	GetRostersIdsByEquipmentId(equipmentId string) ([]string, error)
	DeleteEquipmentInRoster(idRoster string, equipmentId string) error
}

type RosterUpdateRepository interface {
	UpdateRosterStatus(id string, status int32) error
	DeleteUnitInRoster(idRoster string, idUnit string) error
	DeleteEquipmentInRoster(idRoster string, equipmentId string) error
}

type RepositorysProvider interface {
	RosterUpdateRepository() RosterUpdateRepository
}

type RosterUnitOfWork interface {
	RepositorysProvider
	Complete(err error) error
}

type RosterUnitOfWorkFactory interface {
	NewUnitOfWork() (RosterUnitOfWork ,error)
}

type EquipmentRepository interface {
	GetEquipmentInfo(idEquipment string) (Model.EquipmentDetailedInfo, error)
}

type UnitRepository interface {
	GetUnitInfo(idUnit string) (Model.UnitDetailedInfo, error)
}

type RequiredParameters struct {
	Name string
	IdUser string
}

type RosterApp struct {
	db RosterRepository
	unitService UnitRepository
	equipmentService EquipmentRepository
	rosterUnitOfWorkFactory RosterUnitOfWorkFactory
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
	// status roster. 0 - valid, 1 - need update
	Status int32
	Units []UnitAppData
}

type EditRosterAppData struct {
	Name string
	IdUser string
	Units []UnitAppData
}

func CreateRosterApp(db RosterRepository, unitService UnitRepository, equipmentService EquipmentRepository, rosterUnitOfWorkFactory RosterUnitOfWorkFactory) RosterApp {
	return RosterApp{db, unitService, equipmentService, rosterUnitOfWorkFactory}
}

func (app *RosterApp) createRosterAppData(roster Model.Roster) RosterAppData {
	rosterApp := RosterAppData {
		roster.Id,
		roster.Name,
		roster.IdUser,
		roster.Status,
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

func (app *RosterApp) createRosterAppById(id string, status int32, roster EditRosterAppData) RosterAppData {
	rosterApp := RosterAppData {
		id,
		roster.Name,
		roster.IdUser,
		status,
		roster.Units,
	}

	return rosterApp
}

func (app *RosterApp) createRosterInputData(rosterApp RosterAppData) Model.RosterInput {
	rosterIn := Model.RosterInput {
		rosterApp.Id,
		rosterApp.Name,
		rosterApp.IdUser,
		rosterApp.Status,

		[]Model.UnitInput{},
	}

	for i := 0; i < len(rosterApp.Units); i++ {
		unit := rosterApp.Units[i]
		unitApp := Model.UnitInput {
			unit.Id,
			[]Model.EquipmentInput{},
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]
			equipmentApp := Model.EquipmentInput{equipment.Id}
			unitApp.Equipments = append(unitApp.Equipments, equipmentApp)
		}

		rosterIn.Units = append(rosterIn.Units, unitApp)
	}

	return rosterIn
}

func (app *RosterApp) rosterIdExist(id string) bool {
	rosterId, err := app.db.GetRosterIdById(id)
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

// БУДУ ИСПОЛЬЗОВАТЬ ДЛЯ ДАЛЬНЕЙШИХ ПРОВЕРОК НА БИЗНЕС ПРАВИЛА ПРИ ДОБАВЛЕНИИ И ИЗМЕНЕНИИ РОСТЕРА
func (app *RosterApp) getFullRosterInformation(roster Model.Roster) (Model.RosterDetailedInfo, error) {
	rosterDetail := Model.RosterDetailedInfo {
		roster.Id,
		roster.Name,
		roster.IdUser,
		roster.Status,
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

func (app *RosterApp) generateId() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := u.String()
	return id, nil
}

func (app *RosterApp) GetAllRosters() ([]RosterAppData, error) {
	var rosters []RosterAppData
	rostersDB, err := app.db.GetAllRosters()
	if err != nil {
		return rosters, err
	}

	for i := 0; i < len(rostersDB); i++ {
		rosterDB := rostersDB[i]
		rosterApp := app.createRosterAppData(rosterDB)
		rosters = append(rosters, rosterApp)
	}

	return rosters, nil
}

func (app *RosterApp) GetRosterById(id string) (RosterAppData, error) {
	roster, err := app.db.GetRosterById(id)
	if err != nil {
		return RosterAppData {}, err
	}
	rosterApp := app.createRosterAppData(roster)
	return rosterApp, nil
}

func (app *RosterApp) DeleteById(id string) (string, error) {
	err := app.assertRosterNotExist(id)
	if err != nil {
		return "", err
	}
	app.db.Lock()
	deleteId, err := app.db.DeleteRoster(id)
	app.db.Unlock()
	if err != nil {
		return "", err
	}
	return deleteId, nil
}

func (app *RosterApp) assertRosterExists(roster Model.Roster) error {
	rp := RequiredParameters{roster.Name, roster.IdUser}
	id, err:= app.db.GetRosterIdByParams(rp)
	if err != nil {
		return err
	}
	if id != "" {
		return ErrorRosterExist
	}
	return nil
}

func (app *RosterApp) AddNewRoster(editRoster EditRosterAppData) (string, error) {
	id, err := app.generateId()
	if err != nil {
		return "", err
	}
	rosterApp := app.createRosterAppById(id, Model.RosterStatusActive, editRoster)
	rosterInData := app.createRosterInputData(rosterApp)
	roster, err := Model.CreateRoster(rosterInData)
	if err != nil {
		return "", err
	}
	err = app.assertRosterExists(roster)
	if err != nil {
		return "", err
	}
	rosterFullInfo, err:= app.getFullRosterInformation(roster)
	if err != nil {
		return "", err
	}
	err = rosterFullInfo.IsRosterValid()
	if err != nil {
		return "", err
	}
	insertId, err := app.db.InsertNewRoster(roster)
	if err != nil {
		return "", err
	}

	return insertId, nil
}

func (app *RosterApp) assertIsUserCantUpdateRoster(roster RosterAppData) error {
	ownerId, err := app.db.GetOwnerRoster(roster.Id)
	if err != nil {
		return err
	}

	if roster.IdUser != ownerId {
		return ErrorAccessUserRoster
	}

	return nil
}

func (app *RosterApp) UpdateRosterById(id string, editRoster EditRosterAppData) (string, error) {
	err := app.assertRosterNotExist(id)
	if err != nil {
		return "", err
	}
	rosterApp := app.createRosterAppById(id, Model.RosterStatusActive, editRoster)
	err = app.assertIsUserCantUpdateRoster(rosterApp)
	if err != nil {
		return "", err
	}
	rosterInData := app.createRosterInputData(rosterApp)
	roster, err := Model.CreateRoster(rosterInData)
	if err != nil {
		return "", err
	}
	rosterFullInfo, err:= app.getFullRosterInformation(roster)
	if err != nil {
		return "", err
	}
	err = rosterFullInfo.IsRosterValid()
	if err != nil {
		return "", err
	}
	app.db.Lock()
	updateId, err := app.db.UpdateRoster(roster)
	app.db.Unlock()
	if err != nil {
		return "", err
	}

	return updateId, nil
}

func (app *RosterApp) UpdateEquipment(id string) error {
	rostersIds, err := app.db.GetRostersIdsByEquipmentId(id)
	if err != nil {
		return err
	}
	for i := 0; i < len(rostersIds); i++ {
		roster, err := app.db.GetRosterById(rostersIds[i])
		if err != nil {
			return err
		}
		rosterFullInfo, err:= app.getFullRosterInformation(roster)
		if err != nil {
			return err
		}
		err = rosterFullInfo.IsRosterValid()
		if err != nil {
			if err == Model.ErrorRosterCost ||
				err == Model.ErrorCountEquipmernOnUnit ||
				err == Model.ErrorCountEquipmernOnRoster ||
				err == Model.ErrorUnitRole {
				app.db.Lock()
				err = app.db.UpdateRosterStatus(rostersIds[i], Model.RosterStatusNeedUserUpdate)
				app.db.Unlock()
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func (app *RosterApp) DeleteEquipment(id string) error {
	rostersIds, err := app.db.GetRostersIdsByEquipmentId(id)
	if err != nil {
		return err
	}
	for i := 0; i < len(rostersIds); i++ {
		app.db.Lock()
		unitOfWork, err := app.rosterUnitOfWorkFactory.NewUnitOfWork()
		if err != nil {
			app.db.Unlock()
			return err
		}

		rosterUpdate := unitOfWork.RosterUpdateRepository()
		err = rosterUpdate.DeleteEquipmentInRoster(rostersIds[i], id)
		if err != nil {
			unitOfWork.Complete(err)
			app.db.Unlock()
			return err
		}
		err = rosterUpdate.UpdateRosterStatus(rostersIds[i], Model.RosterStatusNeedUserUpdate)
		if err != nil {
			unitOfWork.Complete(err)
			app.db.Unlock()
			return err
		}

		unitOfWork.Complete(nil)
		app.db.Unlock()
	}
	return nil
}

func (app *RosterApp) UpdateUnit(id string) error {
	// TODO СОМНЕВАЮСЬ НУЖНО ЛИ ПОМЕЧАТЬ, ТАК КАК ЛЮБОЕ ОБНОВЛЕНИЕ НЕ ПОВЛИЯЕТ НА ВАЛИДНОСТЬ САМОГО РОСТЕРА
	/*
	rostersIds, err := app.db.GetRostersIdsByUnitId(id)
	if err != nil {
		return err
	}
	for i := 0; i < len(rostersIds); i++ {
		app.mutex.Lock()
		err = app.db.UpdateRosterStatus(rostersIds[i], Model.RosterStatusNeedUserUpdate)
		app.mutex.Unlock()
		if err != nil {
			return err
		}
	}
	*/
	return nil
}

func (app *RosterApp) DeleteUnit(id string) error {
	rostersIds, err := app.db.GetRostersIdsByUnitId(id)
	if err != nil {
		return err
	}
	for i := 0; i < len(rostersIds); i++ {
		app.db.Lock()
		unitOfWork, err := app.rosterUnitOfWorkFactory.NewUnitOfWork()
		if err != nil {
			app.db.Unlock()
			return err
		}

		rosterUpdate := unitOfWork.RosterUpdateRepository()
		err = rosterUpdate.DeleteUnitInRoster(rostersIds[i], id)
		if err != nil {
			unitOfWork.Complete(err)
			app.db.Unlock()
			return err
		}
		err = rosterUpdate.UpdateRosterStatus(rostersIds[i], Model.RosterStatusNeedUserUpdate)
		if err != nil {
			unitOfWork.Complete(err)
			app.db.Unlock()
			return err
		}

		unitOfWork.Complete(nil)
		app.db.Unlock()
	}
	return nil
}