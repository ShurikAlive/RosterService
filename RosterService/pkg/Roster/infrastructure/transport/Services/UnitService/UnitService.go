package UnitService

import (
	Swagger "RosterService/Swagger/UnitService"
	App "RosterService/pkg/Roster/app"
	Model "RosterService/pkg/Roster/model"
	"errors"
	"net/http"
)

var InvalidUnitResponse = errors.New("incorrect unit response")

type UnitService struct {
	UnitApiService *Swagger.APIClient
}

func NewUnitService(unitService *Swagger.APIClient) App.IUnitService {
	return &UnitService{unitService}
}

func (s *UnitService) convertSwaggerUnitToUnitDetailedInfo(unit Swagger.Unit) Model.UnitDetailedInfo {
	unitInfo := Model.UnitDetailedInfo {
		unit.Id,
		unit.Name,
		unit.ForceName,
		unit.Hp,
		unit.Initiative,
		unit.Bs,
		unit.Fs,
		unit.AdditionalRule,
	}

	return unitInfo
}

func (s *UnitService) GetUnitInfo(idUnit string) (Model.UnitDetailedInfo, error) {
	unit, r, err := s.UnitApiService.UnitApi.UnitUnitIdGet( nil, idUnit)
	if err != nil {
		return Model.UnitDetailedInfo{}, err
	}

	if r.StatusCode != http.StatusOK {
		return Model.UnitDetailedInfo{}, InvalidUnitResponse
	}

	return s.convertSwaggerUnitToUnitDetailedInfo(unit), nil
}