package transport

import (
	Swagger "RosterService/Swagger/UnitService"
	"RosterService/cmd/DB"
	App "RosterService/pkg/Roster/app"
	MySQL "RosterService/pkg/Roster/infrastructure/DB"
	"RosterService/pkg/Roster/infrastructure/transport/Services/EquipmentService"
	"RosterService/pkg/Roster/infrastructure/transport/Services/UnitService"
	"github.com/gorilla/mux"
	"net/http"
)

type RosterServer struct {
	app App.RosterApp
}

func CreateRosterServer(connection *DB.Connection, unitService *Swagger.APIClient, equipmentService *Swagger.APIClient) *RosterServer {
	rosterServer := new(RosterServer)
	db := MySQL.NewRosterDB(connection)
	rosterUnitService := UnitService.NewUnitService(unitService)
	rosterEquipmentService := EquipmentService.NewEquipmentService(equipmentService)
	rosterServer.app = App.CreateRosterApp(db, rosterUnitService, rosterEquipmentService)
	return rosterServer
}

func (s *RosterServer) RosterGet(w http.ResponseWriter, r *http.Request) {
	rosters, err := s.app.GetAllRosters()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *RosterServer) RosterPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *RosterServer) RosterRosterIdDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["rosterId"]

	deleteId, err := s.app.DeleteById(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *RosterServer) RosterRosterIdGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["rosterId"]

	roster, err := s.app.GetRosterById(id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *RosterServer) RosterRosterIdPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
