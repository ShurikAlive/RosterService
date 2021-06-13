package transport

import (
	App "RosterService/pkg/roster/app"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RosterEventListener struct {
	channelRebbitMQ *amqp.Channel
	queueRebbitMQ amqp.Queue

	formatter JsonFormatter
	app *App.RosterApp
}

type Event struct {
	Essence string
	TypeEvent string
	IdRecord string
}

func CreateRosterEventListener(channelRebbitMQ *amqp.Channel, queueRebbitMQ amqp.Queue,  app *App.RosterApp) *RosterEventListener {
	formatter := CreateJSONFormatter()
	rosterEventListner := RosterEventListener {channelRebbitMQ, queueRebbitMQ, formatter, app}
	return &rosterEventListner
}

func (listener *RosterEventListener) StartListen() error {
	msgs, err := listener.channelRebbitMQ.Consume(
		listener.queueRebbitMQ.Name, // queue
		"",     // consumer
		false,   // auto-ack
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
				d.Nack(false, true)
				log.Fatal(err)
				continue
			}
			log.Print("Get Event: " + event.Essence + " " + event.TypeEvent + " " + event.IdRecord)

			if event.Essence == "UNIT" && event.TypeEvent == "UPDATE" {
				err = listener.app.UpdateUnit(event.IdRecord)
				if err != nil {
					d.Nack(false, true)
					log.Fatal(err)
					continue
				}
			} else if event.Essence == "UNIT" && event.TypeEvent == "DELETE" {
				err = listener.app.DeleteUnit(event.IdRecord)
				if err != nil {
					d.Nack(false, true)
					log.Fatal(err)
					continue
				}
			} else if event.Essence == "EQUIPMENT" && event.TypeEvent == "UPDATE" {
				err = listener.app.UpdateEquipment(event.IdRecord)
				if err != nil {
					d.Nack(false, true)
					log.Fatal(err)
					continue
				}
			} else if event.Essence == "EQUIPMENT" && event.TypeEvent == "DELETE" {
				err = listener.app.DeleteEquipment(event.IdRecord)
				if err != nil {
					d.Nack(false, true)
					log.Fatal(err)
					continue
				}
			}

			d.Ack(true)
		}
	}()

	return nil
}