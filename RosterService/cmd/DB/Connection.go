package DB

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

type Connection struct {
	Db *sql.DB
}

var ErrorEmptyConnection  = errors.New("Error DB connection")
var ErrorInitConnection  = errors.New("Error Init Connection DB")
var ErrorRecordNotFound = errors.New("Record Not Found")

func (db *Connection) MakeMigrationDB() (error) {
	if db.Db == nil {
		return ErrorEmptyConnection
	}

	driver, err := mysql.WithInstance(db.Db, &mysql.Config{})
	if err != nil {
		return err
	}

	params := ParseEnvDB()

	m, err := migrate.NewWithDatabaseInstance(
		params.DBMigrationsPath,
		params.DBName,
		driver,
	)
	if err != nil {
		return err
	}

	m.Up();

	return nil
}

func ConnectDB() (*sql.DB) {
	params := ParseEnvDB()

	db, err := sql.Open(params.DBType, params.DBUsername + ":" + params.DBPassword + "@/" + params.DBName + "?multiStatements=true")//MySQL80
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil
	}

	return db
}

func InitDB() (*Connection, error) {
	instDB := ConnectDB()
	if instDB == nil {
		return nil, ErrorInitConnection
	}

	var connect = new(Connection)
	connect.Db = instDB
	return connect, nil
}

func (db *Connection) DisconectDB() {
	if db.Db != nil {
		db.Db.Close()
	}
}




