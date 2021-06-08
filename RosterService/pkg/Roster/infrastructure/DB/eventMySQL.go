package DB

import (
	DB "RosterService/pkg/common/infrastructure"
)

type EventSQLDB struct {
	Connection *DB.Connection
}

type EventDB struct {
	IdEvent string
	Essence string
	TypeEvent string
	IdRecord string
}

func CreateEventSQLDB(connection *DB.Connection) *EventSQLDB {
	return &EventSQLDB{connection}
}

func (db *EventSQLDB) InsertEvent(event EventDB) error {
	_, err := db.Connection.Db.Exec("INSERT INTO `roster_db`.`event_execution_tasks` (`idEvent`, `essence`, `typeEvent`, `idRecord`) VALUES (?,?,?,?);", event.IdEvent, event.Essence, event.TypeEvent, event.IdRecord)
	if err != nil {
		return err
	}
	return nil
}

func (db *EventSQLDB) GetEventWithMinimalDate() (EventDB, error) {
	rows, err := db.Connection.Db.Query("SELECT idEvent, essence, typeEvent, idRecord FROM roster_db.event_execution_tasks WHERE creation_time = (SELECT min(creation_time) FROM roster_db.event_execution_tasks) LIMIT 1;")
	if err != nil {
		return EventDB{}, err
	}
	defer rows.Close()

	event :=  EventDB{}
	for rows.Next() {
		err = rows.Scan(
			&event.IdEvent,
			&event.Essence,
			&event.TypeEvent,
			&event.IdRecord)
		if err != nil {
			return EventDB{}, err
		}
	}

	return event, nil
}

func (db *EventSQLDB) DeleteEvent(id string) error {
	_, err := db.Connection.Db.Exec("DELETE FROM `roster_db`.`event_execution_tasks` WHERE idEvent = ?;", id)
	if err != nil {
		return err
	}
	return nil
}


