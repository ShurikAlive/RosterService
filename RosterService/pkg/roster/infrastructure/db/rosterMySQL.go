package DB

import (
	DB "RosterService/pkg/common/infrastructure"
	App "RosterService/pkg/roster/app"
	Model "RosterService/pkg/roster/model"
	"database/sql"
	"github.com/pkg/errors"
)

var lockName = "ROSTER_LOCK"

type MySQLDB struct {
	Connection *DB.Connection
}

type rosterUpdateRepository struct {
	transaction *sql.Tx
}

type rosterUnitOfWork struct{
	transaction *sql.Tx
}

type rosterUnitOfWorkFactory struct{
	Connection *DB.Connection
}

func NewRosterDB(connection *DB.Connection) App.RosterRepository {
	if connection.Db == nil {
		return nil
	}
	return &MySQLDB{connection}
}

func NewRosterUnitOfWorkFactory(connection *DB.Connection) App.RosterUnitOfWorkFactory {
	if connection.Db == nil {
		return nil
	}
	return &rosterUnitOfWorkFactory{connection}
}

func (u *rosterUnitOfWork) RosterUpdateRepository() App.RosterUpdateRepository {
	return &rosterUpdateRepository{u.transaction}
}

func (u *rosterUnitOfWork) Complete(err error) error {
	if err != nil{
		err2 := u.transaction.Rollback()
		if err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
		return err
	}
	return errors.WithStack(u.transaction.Commit())
}

func (uf *rosterUnitOfWorkFactory) NewUnitOfWork() (App.RosterUnitOfWork ,error) {
	transaction, err := uf.Connection.Db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &rosterUnitOfWork{transaction}, nil
}


func (db *MySQLDB) Lock() error {
	_, err := db.Connection.Db.Exec("SELECT GET_LOCK(?, -1)", lockName)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLDB) Unlock() error {
	_, err := db.Connection.Db.Exec("SELECT RELEASE_LOCK(?)", lockName)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLDB) getUnitEquipments(rosterId string, unitId string) ([]Model.Equipment, error) {
	var equipments []Model.Equipment

	rowsEquipments, err := db.Connection.Db.Query("SELECT idEquipment FROM roster_db.roster_equipments WHERE idRoster = ? AND idUnit = ?;", rosterId, unitId)
	if err != nil {
		return []Model.Equipment{}, err
	}
	defer rowsEquipments.Close()

	for rowsEquipments.Next() {
		equipment := Model.Equipment{}
		err = rowsEquipments.Scan(&equipment.Id)
		if err != nil {
			return []Model.Equipment{}, err
		}
		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

func (db *MySQLDB) getRosterUnits(rosterId string) ([]Model.Unit, error) {
	var units []Model.Unit

	rowsUnits, err := db.Connection.Db.Query("SELECT idUnit FROM roster_db.roster_units WHERE idRoster = ?;", rosterId)
	if err != nil {
		return []Model.Unit{}, err
	}
	defer rowsUnits.Close()

	for rowsUnits.Next() {
		unit := Model.Unit{}
		err = rowsUnits.Scan(&unit.Id)
		if err != nil {
			return []Model.Unit{}, err
		}

		equipments, err := db.getUnitEquipments(rosterId, unit.Id)
		if err != nil {
			return []Model.Unit{}, err
		}
		unit.Equipments = equipments

		units = append(units, unit)
	}

	return units, nil
}

func (db *MySQLDB) GetAllRosters() ([]Model.Roster, error) {
	rows, err := db.Connection.Db.Query("SELECT * FROM roster_db.rosters;")
	if err != nil {
		return []Model.Roster{}, err
	}
	defer rows.Close()

	var rosters []Model.Roster

	for rows.Next() {
		roster := Model.Roster{}
		err = rows.Scan(&roster.Id,
			&roster.Name,
			&roster.IdUser,
			&roster.Status)
		if err != nil {
			return []Model.Roster{}, err
		}

		units, err := db.getRosterUnits(roster.Id)
		if err != nil {
			return []Model.Roster{}, err
		}
		roster.Units = units

		rosters = append(rosters, roster)
	}

	return rosters, nil
}

func (db *MySQLDB) GetRosterById(id string) (Model.Roster, error) {
	rows, err := db.Connection.Db.Query("SELECT * FROM roster_db.rosters WHERE id = ?;", id)
	if err != nil {
		return Model.Roster{}, err
	}
	defer rows.Close()

	roster := Model.Roster{}
	for rows.Next() {
		err = rows.Scan(&roster.Id,
			&roster.Name,
			&roster.IdUser,
			&roster.Status)
		if err != nil {
			return Model.Roster{}, err
		}
	}
	units, err := db.getRosterUnits(roster.Id)
	if err != nil {
		return Model.Roster{}, err
	}
	roster.Units = units

	return roster, nil
}

func (db *MySQLDB) DeleteRoster(id string) (string, error) {
	tx, err := db.Connection.Db.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`roster_equipments` WHERE idRoster = ?;", id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`roster_units` WHERE idRoster = ?;", id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`rosters` WHERE id = ?;", id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return id, nil
}

func (db *MySQLDB) GetRosterIdById(id string) (string, error) {
	rows, err := db.Connection.Db.Query("SELECT id FROM roster_db.rosters WHERE id = ?;", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	idDB := ""
	for rows.Next() {
		err = rows.Scan(&idDB)
		if err != nil {
			return "", err
		}
	}

	return idDB, nil
}

func (db *MySQLDB) GetRosterIdByParams(params App.RequiredParameters) (string, error) {
	rows, err := db.Connection.Db.Query("SELECT id FROM roster_db.rosters WHERE idUser = ? AND Name = ?;", params.IdUser, params.Name)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	idDB := ""
	for rows.Next() {
		err = rows.Scan(&idDB)
		if err != nil {
			return "", err
		}
	}

	return idDB, nil
}

func (db *MySQLDB) InsertNewRoster(roster Model.Roster) (string, error) {
	tx, err := db.Connection.Db.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec("INSERT INTO `roster_db`.`rosters` (`id`,`Name`,	`idUser`, `status`) VALUES (?,?,?,?);", roster.Id, roster.Name, roster.IdUser, roster.Status)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	for i := 0; i < len(roster.Units); i++ {
		unit := roster.Units[i]

		_, err = tx.Exec("INSERT INTO `roster_db`.`roster_units` (`idRoster`,`idUnit`) VALUES (?,?);", roster.Id, unit.Id)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]

			_, err = tx.Exec("INSERT INTO `roster_db`.`roster_equipments` (`idRoster`,`idUnit`,`idEquipment`) VALUES (?,?,?);",
				roster.Id, unit.Id, equipment.Id)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return roster.Id, nil
}

func (db *MySQLDB) UpdateRoster(roster Model.Roster) (string, error) {
	tx, err := db.Connection.Db.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`roster_equipments` WHERE idRoster = ?;", roster.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`roster_units` WHERE idRoster = ?;", roster.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec("UPDATE `roster_db`.`rosters` SET `Name` = ?, `idUser` = ?, `status` = ? WHERE `id` = ?;", roster.Name, roster.IdUser, roster.Status, roster.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	for i := 0; i < len(roster.Units); i++ {
		unit := roster.Units[i]

		_, err = tx.Exec("INSERT INTO `roster_db`.`roster_units` (`idRoster`,`idUnit`) VALUES (?,?);", roster.Id, unit.Id)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		for j := 0; j < len(unit.Equipments); j++ {
			equipment := unit.Equipments[j]

			_, err = tx.Exec("INSERT INTO `roster_db`.`roster_equipments` (`idRoster`,`idUnit`,`idEquipment`) VALUES (?,?,?);",
				roster.Id, unit.Id, equipment.Id)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return roster.Id, nil
}

func (db *MySQLDB) GetOwnerRoster(id string) (string, error) {
	rows, err := db.Connection.Db.Query("SELECT idUser FROM roster_db.rosters WHERE id = ?;", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	userId := ""
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return "", err
		}
	}

	return userId, nil
}

func (db *MySQLDB) GetRostersIdsByUnitId(unitId string) ([]string, error) {
	rows, err := db.Connection.Db.Query("SELECT DISTINCT idRoster FROM roster_db.roster_units WHERE idUnit = ?;", unitId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

func (db *MySQLDB) UpdateRosterStatus(id string, status int32) error {
	_, err := db.Connection.Db.Exec("UPDATE `roster_db`.`rosters` SET `status` = ? WHERE `id` = ?;", status, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLDB) DeleteUnitInRoster(idRoster string, idUnit string) error {
	tx, err := db.Connection.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`roster_units` WHERE idRoster = ? AND idUnit = ?;", idRoster, idUnit)
	if err != nil {
		tx.Rollback()
		return  err
	}

	_, err = tx.Exec("DELETE FROM `roster_db`.`roster_equipments` WHERE idRoster = ? AND idUnit = ?;", idRoster, idUnit)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (db *MySQLDB) GetRostersIdsByEquipmentId(equipmentId string) ([]string, error) {
	rows, err := db.Connection.Db.Query("SELECT DISTINCT idRoster FROM roster_db.roster_equipments WHERE idEquipment = ?;", equipmentId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

func (db *MySQLDB) DeleteEquipmentInRoster(idRoster string, equipmentId string) error {
	_, err := db.Connection.Db.Exec("DELETE FROM `roster_db`.`roster_equipments` WHERE idRoster = ? AND idEquipment = ?;", idRoster, equipmentId)
	if err != nil {
		return  err
	}
	return nil
}

func (db *rosterUpdateRepository) UpdateRosterStatus(id string, status int32) error {
	_, err := db.transaction.Exec("UPDATE `roster_db`.`rosters` SET `status` = ? WHERE `id` = ?;", status, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *rosterUpdateRepository) DeleteUnitInRoster(idRoster string, idUnit string) error {

	_, err := db.transaction.Exec("DELETE FROM `roster_db`.`roster_units` WHERE idRoster = ? AND idUnit = ?;", idRoster, idUnit)
	if err != nil {
		return  err
	}

	_, err = db.transaction.Exec("DELETE FROM `roster_db`.`roster_equipments` WHERE idRoster = ? AND idUnit = ?;", idRoster, idUnit)
	if err != nil {
		return err
	}

	return nil
}

func (db *rosterUpdateRepository) DeleteEquipmentInRoster(idRoster string, equipmentId string) error {
	_, err := db.transaction.Exec("DELETE FROM `roster_db`.`roster_equipments` WHERE idRoster = ? AND idEquipment = ?;", idRoster, equipmentId)
	if err != nil {
		return  err
	}

	return nil
}