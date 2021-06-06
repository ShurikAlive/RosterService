package infrastructure

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

type Connection struct {
	Db *sql.DB

	dbType string
	dbUsername string
	dbPassword string
	dbName string
}

func (db *Connection) MakeMigrationDB(migrationPath string) (error) {
	driver, err := mysql.WithInstance(db.Db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		db.dbName,
		driver,
	)
	if err != nil {
		return err
	}

	m.Up()

	return nil
}

func connectDB(dbType string, dbUsername string, dbPassword string, dbName string) (*sql.DB, error) {

	db, err := sql.Open(dbType, dbUsername + ":" + dbPassword + "@/" + dbName + "?multiStatements=true")//MySQL80
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}

func InitDB(dbType string, dbUsername string, dbPassword string, dbName string) (*Connection, error) {
	instDB, err := connectDB(dbType, dbUsername, dbPassword, dbName)
	if err != nil {
		return nil, err
	}

	var connect = new(Connection)
	connect.Db = instDB
	connect.dbType = dbType
	connect.dbUsername = dbUsername
	connect.dbPassword = dbPassword
	connect.dbName = dbName
	return connect, nil
}

func (db *Connection) DisconectDB() {
	if db.Db != nil {
		db.Db.Close()
	}
}




