package main

import (
	DB "RosterService/pkg/common/infrastructure"
	Roster "RosterService/pkg/roster/infrastructure/transport"
	Swagger "RosterService/swagger/go"
	SwaggerUnitService "RosterService/swagger/unitService"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverParameters := initServerParameters()
	initLogFile()

	db, err := DB.InitDB(serverParameters.DBType, serverParameters.DBUsername, serverParameters.DBPassword, serverParameters.DBName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.DisconectDB()

	err = db.MakeMigrationDB(serverParameters.DBMigrationsPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := amqp.Dial(serverParameters.ConnectRabbitMQ)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		serverParameters.QueueNameRabbitMQ, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	killSignalChan := getKillSignalChan()
	srv, err := startServer(db, serverParameters, ch, q)
	if err != nil {
		log.Fatal(err)
		return
	}
	waitForKillSignal(killSignalChan)
	srv.Shutdown(context.Background())
}


func initServerParameters() (*Config) {
	serverParameters, err := ParseEnv()
	if err != nil {
		log.Fatal("Cannot init server parameters!")
	}
	return serverParameters
}

func initLogFile() {
	log.SetFormatter(&log.JSONFormatter{})
}

func InitUnitService(conf *Config) *SwaggerUnitService.APIClient {
	config := &SwaggerUnitService.Configuration{
		BasePath:      conf.UnitServiceBasePath,
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
	}
	return SwaggerUnitService.NewAPIClient(config)
}

func InitEquipmentService(conf *Config) *SwaggerUnitService.APIClient {
	config := &SwaggerUnitService.Configuration{
		BasePath:      conf.EquipmentServiceBasePath,
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
	}
	return SwaggerUnitService.NewAPIClient(config)
}

func InitRosterHendlerFunc(router *mux.Router, connection *DB.Connection, conf *Config, channelRebbitMQ *amqp.Channel, queueRebbitMQ amqp.Queue) (*mux.Router, error) {
	unitService := InitUnitService(conf)
	equipmentService := InitEquipmentService(conf)
	rosterServer, err := Roster.CreateRosterServer(connection, unitService, equipmentService, channelRebbitMQ, queueRebbitMQ)
	if err != nil {
		return nil, err
	}

	unitHandlerFuncs := map[string]http.HandlerFunc {
		"RosterGet" : rosterServer.RosterGet,
		"RosterPost" : rosterServer.RosterPost,
		"RosterRosterIdDelete" : rosterServer.RosterIdDelete,
		"RosterRosterIdGet" : rosterServer.RosterIdGet,
		"RosterRosterIdPut" : rosterServer.RosterIdPut,
	}

	for name, unitHendlerFunc := range unitHandlerFuncs {
		router.GetRoute(name).Handler(unitHendlerFunc)
	}

	return router, nil
}

func startServer(connection *DB.Connection, conf *Config, channelRebbitMQ *amqp.Channel, queueRebbitMQ amqp.Queue) (*http.Server, error) {
	router := Swagger.NewRouter()
	router, err := InitRosterHendlerFunc(router, connection, conf, channelRebbitMQ, queueRebbitMQ)
	if err != nil {
		return nil, err
	}
	log.Fatal(http.ListenAndServe(conf.ServeRESTAddress, router))
	srv := &http.Server{Addr: conf.ServeRESTAddress, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv, nil
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}