package model

import "errors"

var ErrorRosterCost = errors.New("error roster cost!")
var ErrorCountEquipmernOnUnit = errors.New("error count equipment on unit!")
var ErrorCountEquipmernOnRoster = errors.New("error count equipment on roster!")
var ErrorUnitRole = errors.New("error unit role in rostern!")

var MAX_COST_ROSTER = 12

type UnitDetailedInfo struct {
	Id string
	Name string
	ForceName string
	// count heals point unit
	Hp int32
	// initiative unit
	Initiative int32
	// ability to shoot unit
	Bs int32
	// ability to fight unit
	Fs int32
	// Additionat ability soldes
	AdditionalRule string
}

type EquipmentDetailedInfo struct {
	Id string
	Name string
	// limit equipment on one unit. -1 - unlimit
	LimitOnUnit int32
	// limit equipment on one team. -1 - unlimit
	LimitOnTeam int32
	// The role of a soldier available when selecting ammunition.
	SoldarRole string
	// game rule equipment
	Rule string
	// limit equipment on game. -1 - unlimit
	Ammo int32
	// cost equipment in game points
	Cost int32
}

type RosterUnitDetailedInfo struct {
	UnitInfo UnitDetailedInfo
	EquipmentsInfo []EquipmentDetailedInfo
}

type RosterDetailedInfo struct {
	Id string
	Name string
	IdUser string
	// status roster. 0 - valid, 1 - need update
	Status int32
	Units []RosterUnitDetailedInfo
}


func (r *RosterDetailedInfo) costRosterValid() error {
	costRoster := int32(0)
	for i := 0; i < len(r.Units); i++ {
		unit := r.Units[i]
		for j := 0; j < len(unit.EquipmentsInfo); j++ {
			costRoster += unit.EquipmentsInfo[j].Cost
		}
	}

	if costRoster > int32(MAX_COST_ROSTER) {
		return ErrorRosterCost
	}

	return nil
}

func (r *RosterDetailedInfo) countEquipmentOnUnitRosterValid() error {
	for i := 0; i < len(r.Units); i++ {
		unit := r.Units[i]
		unitEqipmentCount := make(map[string]int32)
		for j := 0; j < len(unit.EquipmentsInfo); j++ {
			eqipment := unit.EquipmentsInfo[j]
			if eqipment.LimitOnUnit != -1 {
				_, equipmentExistInMap := unitEqipmentCount[eqipment.Id]
				if equipmentExistInMap {
					unitEqipmentCount[eqipment.Id] -= 1
					if unitEqipmentCount[eqipment.Id] < 0 {
						return ErrorCountEquipmernOnUnit
					}
				} else {
					unitEqipmentCount[eqipment.Id] = eqipment.LimitOnUnit - 1
				}
			}
		}
	}

	return nil
}

func (r *RosterDetailedInfo) countEquipmentOnRosterValid() error {
	rosterEqipmentCount := make(map[string]int32)
	for i := 0; i < len(r.Units); i++ {
		unit := r.Units[i]
		for j := 0; j < len(unit.EquipmentsInfo); j++ {
			equipment := unit.EquipmentsInfo[j]
			if equipment.LimitOnTeam != -1 {
				_, equipmentExistInMap := rosterEqipmentCount[equipment.Id]
				if equipmentExistInMap {
					rosterEqipmentCount[equipment.Id] -= 1
					if rosterEqipmentCount[equipment.Id] < 0 {
						return ErrorCountEquipmernOnRoster
					}
				} else {
					rosterEqipmentCount[equipment.Id] = equipment.LimitOnTeam - 1
				}
			}
		}
	}
	return nil
}

func (r *RosterDetailedInfo) unitRolesInRosterValid() error {
	for i := 0; i < len(r.Units); i++ {
		unit := r.Units[i]
		unitRole := ""
		for j := 0; j < len(unit.EquipmentsInfo); j++ {
			equipment := unit.EquipmentsInfo[j]
			if equipment.SoldarRole != "" {
				if unitRole == "" {
					unitRole = equipment.SoldarRole
				} else {
					return ErrorUnitRole
				}
			}
		}
	}
	return nil
}

func (r *RosterDetailedInfo) IsRosterValid() error {

	err := r.costRosterValid()
	if err != nil {
		return err
	}

	err = r.countEquipmentOnUnitRosterValid()
	if err != nil {
		return err
	}

	err = r.countEquipmentOnRosterValid()
	if err != nil {
		return err
	}

	err = r.unitRolesInRosterValid()
	if err != nil {
		return err
	}

	return nil
}