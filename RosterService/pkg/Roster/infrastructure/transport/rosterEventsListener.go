package transport

import (
	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	DB "RosterService/pkg/roster/infrastructure/db"
)

type RosterEventListener struct {
	channelRebbitMQ *amqp.Channel
	queueRebbitMQ amqp.Queue

	formatter JsonFormatter
	db *DB.EventSQLDB
}

type Event struct {
	Essence string
	TypeEvent string
	IdRecord string
}

func CreateRosterEventListener(channelRebbitMQ *amqp.Channel, queueRebbitMQ amqp.Queue, db *DB.EventSQLDB) *RosterEventListener {
	formatter := CreateJSONFormatter()
	rosterEventListner := RosterEventListener {channelRebbitMQ, queueRebbitMQ, formatter, db}
	return &rosterEventListner
}

func (listener *RosterEventListener) generateId() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := u.String()
	return id, nil
}


func (listener *RosterEventListener) createEventDB(id string, event Event) DB.EventDB {
	eventDB := DB.EventDB {
		id,
		event.Essence,
		event.TypeEvent,
		event.IdRecord,
	}
	return eventDB
}

func (listener *RosterEventListener) StartListen() error {
	msgs, err := listener.channelRebbitMQ.Consume(
		listener.queueRebbitMQ.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			event, err := listener.formatter.ConvertEvent(d.Body)
			log.Print(d.Body)
			if err != nil {
				log.Fatal(err)
				continue
			}
			log.Print("Get Event: " + event.Essence + " " + event.TypeEvent + " " + event.IdRecord)
			id, err := listener.generateId()
			if err != nil {
				log.Fatal(err)
				continue
			}
			eventDB := DB.EventDB {
				id,
				event.Essence,
				event.TypeEvent,
				event.IdRecord,
			}
			err = listener.db.InsertEvent(eventDB)
			if err != nil {
				log.Fatal(err)
				continue
			}
			log.Print("Event inserted in DB: " + eventDB.IdRecord  + " " + event.Essence + " " + event.TypeEvent + " " + event.IdRecord)
		}
	}()

	return nil
}