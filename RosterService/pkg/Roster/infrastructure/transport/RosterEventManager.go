package transport

import (
	App "RosterService/pkg/roster/app"
	DB "RosterService/pkg/roster/infrastructure/db"
	Connection "RosterService/pkg/common/infrastructure"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RosterEventManager struct {
	app *App.RosterApp
	listener *RosterEventListener
	db *DB.EventSQLDB
}

func CreateRosterEventManager(channelRebbitMQ *amqp.Channel, queueRebbitMQ amqp.Queue, app *App.RosterApp, conDB *Connection.Connection) RosterEventManager {
	db := DB.CreateEventSQLDB(conDB)
	eventListener := CreateRosterEventListener(channelRebbitMQ, queueRebbitMQ, db)
	err := eventListener.StartListen()
	if err != nil {
		return RosterEventManager{}
	}
	eventManager := RosterEventManager{app, eventListener, db}
	return eventManager
}

func (manager *RosterEventManager) StartEventHandling() {
	go func() {
		for manager.db.Connection.Db != nil {
			event, err := manager.db.GetEventWithMinimalDate()
			if err != nil {
				log.Fatal(err)
				continue
			}
			if event.IdEvent == "" {
				continue
			}

			log.Print("Execute Event: " + event.IdEvent + " " + event.Essence + " " + event.TypeEvent + " " + event.IdRecord)

			if event.Essence == "UNIT" && event.TypeEvent == "UPDATE" {
				err = manager.app.UpdateUnit(event.IdRecord)
				if err != nil {
					log.Fatal(err)
					continue
				}
			} else if event.Essence == "UNIT" && event.TypeEvent == "DELETE" {
				err = manager.app.DeleteUnit(event.IdRecord)
				if err != nil {
					log.Fatal(err)
					continue
				}
			} else if event.Essence == "EQUIPMENT" && event.TypeEvent == "UPDATE" {
				err = manager.app.UpdateEquipment(event.IdRecord)
				if err != nil {
					log.Fatal(err)
					continue
				}
			} else if event.Essence == "EQUIPMENT" && event.TypeEvent == "DELETE" {
				err = manager.app.DeleteEquipment(event.IdRecord)
				if err != nil {
					log.Fatal(err)
					continue
				}
			}

			err = manager.db.DeleteEvent(event.IdEvent)
			if err != nil {
				log.Fatal(err)
				continue
			}
		}
	}()
}
/*


 */