package model

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
	Units []RosterUnitDetailedInfo
}