package transport

import (
	DB "RosterService/pkg/common/infrastructure"
	App "RosterService/pkg/roster/app"
	MySQL "RosterService/pkg/roster/infrastructure/db"
	"RosterService/pkg/roster/infrastructure/transport/services/equipmentService"
	"RosterService/pkg/roster/infrastructure/transport/services/unitService"
	Model "RosterService/pkg/roster/model"
	Swagger "RosterService/swagger/unitService"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type RosterServer struct {
	app App.RosterApp
	formatter JsonFormatter
}

func CreateRosterServer(connection *DB.Connection, unitService *Swagger.APIClient, equipmentService *Swagger.APIClient) *RosterServer {
	rosterServer := new(RosterServer)
	db := MySQL.NewRosterDB(connection)
	rosterUnitService := UnitService.NewUnitService(unitService)
	rosterEquipmentService := EquipmentService.NewEquipmentService(equipmentService)
	rosterServer.app = App.CreateRosterApp(db, rosterUnitService, rosterEquipmentService)
	rosterServer.formatter = CreateJSONFormatter()
	return rosterServer
}

func (s *RosterServer) getErrorCode(err error) int {
	code := http.StatusInternalServerError
	switch err {
	case App.ErrorRosterExist:
		code = http.StatusBadRequest
	case App.ErrorRosterIdNotExist:
		code = http.StatusBadRequest

	case Model.ErrorCountUnitInRoster:
		code = http.StatusBadRequest
	case Model.ErrorRosterUserId:
		code = http.StatusBadRequest
	case Model.ErrorCountEquipmernOnRoster:
		code = http.StatusBadRequest
	case Model.ErrorCountEquipmernOnUnit:
		code = http.StatusBadRequest
	case Model.ErrorDoubleUnit:
		code = http.StatusBadRequest
	case Model.ErrorRosterCost:
		code = http.StatusBadRequest
	case Model.ErrorRosterName:
		code = http.StatusBadRequest
	case Model.ErrorUnitRole:
		code = http.StatusBadRequest
	}

	return code
}

func (s *RosterServer) RosterGet(w http.ResponseWriter, r *http.Request) {
	rosters, err := s.app.GetAllRosters()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), s.getErrorCode(err))
		return
	}

	json, err := s.formatter.ConvertAllRostersAppDataToJSON(rosters)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, json)
}

func (s *RosterServer) RosterPost(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	rosterEdit, err := s.formatter.ConvertJsonToRocterEditAppData(b)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := s.app.AddNewRoster(rosterEdit)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), s.getErrorCode(err))
		return
	}

	idJSON:= s.formatter.ConvertIdToJSON(id)

	fmt.Fprintf(w,idJSON)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *RosterServer) RosterIdDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["rosterId"]

	deleteId, err := s.app.DeleteById(id)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), s.getErrorCode(err))
		return
	}

	idJSON:= s.formatter.ConvertIdToJSON(deleteId)

	fmt.Fprintf(w,idJSON)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *RosterServer) RosterIdGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["rosterId"]

	roster, err := s.app.GetRosterById(id)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := s.formatter.ConvertRosterAppDataToJSON(roster)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSONResponse(w, json)
}

func (s *RosterServer) RosterIdPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["rosterId"]

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	rosterEdit, err := s.formatter.ConvertJsonToRocterEditAppData(b)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idEdit, err := s.app.UpdateRosterById(id, rosterEdit)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), s.getErrorCode(err))
		return
	}

	idJSON:= s.formatter.ConvertIdToJSON(idEdit)

	fmt.Fprintf(w,idJSON)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func JSONResponse(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, string(json))
	if err != nil {
		log.WithField("err", err).Error("write response error")
	}
}