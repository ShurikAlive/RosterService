package EquipmentService

import (
	App "RosterService/pkg/roster/app"
	Model "RosterService/pkg/roster/model"
	Swagger "RosterService/swagger/unitService"
	"errors"
	"net/http"
)

var InvalidEquipmentResponse = errors.New("incorrect equipment response")

type EquipmentService struct {
	EquipmentApiService *Swagger.APIClient
}

func NewEquipmentService(equipmentService *Swagger.APIClient) App.EquipmentRepository {
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