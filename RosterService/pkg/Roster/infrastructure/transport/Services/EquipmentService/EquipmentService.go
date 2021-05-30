package EquipmentService

import (
	Swagger "RosterService/Swagger/UnitService"
	App "RosterService/pkg/Roster/app"
	Model "RosterService/pkg/Roster/model"
	"errors"
	"net/http"
)

var InvalidEquipmentResponse = errors.New("incorrect equipment response")

type EquipmentService struct {
	EquipmentApiService *Swagger.APIClient
}

func NewEquipmentService(equipmentService *Swagger.APIClient) App.IEquipmentService {
	return &EquipmentService{equipmentService}
}

func (s *EquipmentService) convertSwaggerEquipmentToEquipmentDetailedInfo(equipment Swagger.Equipment) Model.EquipmentDetailedInfo {
	equipmentInfo := Model.EquipmentDetailedInfo {
		equipment.Id,
		equipment.Name,
		equipment.LimitOnUnit,
		equipment.LimitOnTeam,
		equipment.SoldarRole,
		equipment.Rule,
		equipment.Ammo,
		equipment.Cost,
	}

	return equipmentInfo
}

func (s *EquipmentService) GetEquipmentInfo(idEquipment string) (Model.EquipmentDetailedInfo, error) {
	equipment, r, err := s.EquipmentApiService.WargearsApi.EquipmentEquipmentIdGet( nil, idEquipment)
	if err != nil {
		return Model.EquipmentDetailedInfo{}, err
	}

	if r.StatusCode != http.StatusOK {
		return Model.EquipmentDetailedInfo{}, InvalidEquipmentResponse
	}

	return s.convertSwaggerEquipmentToEquipmentDetailedInfo(equipment), nil
}