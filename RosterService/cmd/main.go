package cmd

import (
	"RosterService/cmd/DB"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	Swagger "RosterService/Swagger/go"
	SwaggerUnitService "RosterService/Swagger/UnitService"
	Roster "RosterService/pkg/Roster/infrastructure/transport"
)

func main() {
	serverParameters := initServerParameters()
	initLogFile()

	db, err := DB.InitDB()
	if err != nil {
		return
	}
	defer db.DisconectDB()

	killSignalChan := getKillSignalChan()
	srv := startServer(db, serverParameters)
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
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}
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

func InitRosterHendlerFunc(router *mux.Router, connection *DB.Connection, conf *Config) (*mux.Router) {
	unitService := InitUnitService(conf)
	equipmentService := InitEquipmentService(conf)
	rosterServer := Roster.CreateRosterServer(connection, unitService, equipmentService)

	unitHandlerFuncs := map[string]http.HandlerFunc {
		"RosterGet" : rosterServer.RosterGet,
		"RosterPost" : rosterServer.RosterPost,
		"RosterRosterIdDelete" : rosterServer.RosterRosterIdDelete,
		"RosterRosterIdGet" : rosterServer.RosterRosterIdGet,
		"RosterRosterIdPut" : rosterServer.RosterRosterIdPut,
	}

	for name, unitHendlerFunc := range unitHandlerFuncs {
		router.GetRoute(name).Handler(unitHendlerFunc)
	}

	return router
}

func startServer(connection *DB.Connection, conf *Config) *http.Server {
	router := Swagger.NewRouter()
	router = InitRosterHendlerFunc(router, connection, conf)
	log.Fatal(http.ListenAndServe(conf.ServeRESTAddress, router))
	srv := &http.Server{Addr: conf.ServeRESTAddress, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
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