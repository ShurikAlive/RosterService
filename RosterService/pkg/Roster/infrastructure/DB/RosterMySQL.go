package DB

import (
	"RosterService/cmd/DB"
	App "RosterService/pkg/Roster/app"
	Model "RosterService/pkg/Roster/model"
)

type MySQLDB struct {
	Connection *DB.Connection
}

func NewRosterDB(connection *DB.Connection) App.IRosterDB {
	if connection.Db == nil {
		return nil
	}
	return &MySQLDB{connection}
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
			&roster.IdUser)
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

func (db *MySQLDB) GetRosterInDBById(id string) (Model.Roster, error) {
	rows, err := db.Connection.Db.Query("SELECT * FROM roster_db.rosters WHERE id = ?;", id)
	if err != nil {
		return Model.Roster{}, err
	}
	defer rows.Close()

	roster := Model.Roster{}
	for rows.Next() {
		roster := Model.Roster{}
		err = rows.Scan(&roster.Id,
			&roster.Name,
			&roster.IdUser)
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

func (db *MySQLDB) GetRosterIdInDBById(id string) (string, error) {
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