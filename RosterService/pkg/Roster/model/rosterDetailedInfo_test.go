package model

import "testing"

func createValidRosterDetailedInfo() RosterDetailedInfo {
	roster := RosterDetailedInfo{
		"test roster id",
		"test roster",
		"test user id",
		0,
		[]RosterUnitDetailedInfo{
			{
				UnitDetailedInfo {
					"4c2732e8-522c-4cda-80b2-2a7762eb4ee2",
					"Robert A. Pepper",
					"Special Forces Airborne",
					4,
					8,
					1,
					4,
					"",
				},
				[]EquipmentDetailedInfo{
					{
						"780fa83b-dacf-4d84-b0a8-13620ed57b35",
						"M230",
						1,
						1,
						"Grenader",
						"2 orders: The under barrel grenade launcher deals 1 wound to all models in its AoE. The target must be in LoS and does not have to be a model. No max range to this tossed weapon.",
						2,
						3,
					},
				},
			},
			{
				UnitDetailedInfo {
					"83b1a1e7-0927-459c-905f-57789f9cf2ec",
					"Robert R. River",
					"Special Forces Airborne",
					4,
					8,
					3,
					2,
					"",
				},
				[]EquipmentDetailedInfo{
					{
						"5131d405-88fc-4ee8-aa1b-f92da9dbe883",
						"Machine gunner M249 LMG",
						1,
						1,
						"Machine gunner",
						"This soldier is the fireteam is gunner. The M249 replaces the M4 carbine as his primary weapon. This soldier can perform Suppression Fire for 1 order. This weapon deals 2 wounds. This soldier receives a -2 to his accuracy.",
						-1,
						3,
					},
				},
			},
			{
				UnitDetailedInfo {
					"780b1856-951b-4ed6-a93d-4a3cd16a93b9",
					"Patrick Singletary",
					"Special Forces Airborne",
					4,
					6,
					2,
					2,
					"",
				},
				[]EquipmentDetailedInfo{
					{
						"b1256ae9-b91e-4ef0-8283-87c1687a6c4b",
						"M67 Grenade",
						-1,
						-1,
						"",
						"1 order: This tossed weapon deals 1 wounds to all models in its AoE.",
						1,
						2,
					},
				},
			},
			{
				UnitDetailedInfo {
					"aa6f736f-4dfd-4e06-b16c-3bccd2d99fd9",
					"Gregiory W. Morris",
					"Special Forces Airborne",
					4,
					8,
					1,
					1,
					"When Morris moves into hand to hand combat his roll is at +2",
				},
				[]EquipmentDetailedInfo{
					{
						"0eddeb73-801d-4e58-b795-a1334d5d2879",
						"M84 Stun Grenade",
						-1,
						-1,
						"",
						"1 orders: This tossed weapon renders all models in its AoE inactive",
						2,
						2,
					},
					{
						"b1256ae9-b91e-4ef0-8283-87c1687a6c4b",
						"M67 Grenade",
						-1,
						-1,
						"",
						"1 order: This tossed weapon deals 1 wounds to all models in its AoE.",
						1,
						2,
					},
				},
			},
		},
	}

	return roster
}

func TestCheckValidationFuncRoster(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	err := roster.IsRosterValid()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidCostRoster(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	granade := EquipmentDetailedInfo {
		"b1256ae9-b91e-4ef0-8283-87c1687a6c4b",
		"M67 Grenade",
		-1,
		-1,
		"",
		"1 order: This tossed weapon deals 1 wounds to all models in its AoE.",
		1,
		2,
	}
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, granade)
	err := roster.IsRosterValid()
	if err != ErrorRosterCost {
		t.Fatal(err)
	}
}

func TestInvalidCountEquipmentsOnUnit(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	testEquipment := EquipmentDetailedInfo {
		"test id",
		"Test equipment",
		2,
		-1,
		"",
		"Test rule",
		-1,
		0,
	}
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)

	err := roster.IsRosterValid()
	if err != ErrorCountEquipmernOnUnit {
		t.Fatal(err)
	}
}

func TestValidCountEquipmentsOnUnit(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	testEquipment := EquipmentDetailedInfo {
		"test id",
		"Test equipment",
		2,
		-1,
		"",
		"Test rule",
		-1,
		0,
	}
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)

	err := roster.IsRosterValid()
	if err == ErrorCountEquipmernOnUnit {
		t.Fatal(err)
	}
}


func TestInvalidCountEquipmentsOnRoster(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	testEquipment := EquipmentDetailedInfo {
		"test id",
		"Test equipment",
		-1,
		2,
		"",
		"Test rule",
		-1,
		0,
	}
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)
	roster.Units[2].EquipmentsInfo = append(roster.Units[2].EquipmentsInfo, testEquipment)
	roster.Units[3].EquipmentsInfo = append(roster.Units[3].EquipmentsInfo, testEquipment)

	err := roster.IsRosterValid()
	if err != ErrorCountEquipmernOnRoster {
		t.Fatal(err)
	}
}

func TestValidCountEquipmentsOnRoster(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	testEquipment := EquipmentDetailedInfo {
		"test id",
		"Test equipment",
		-1,
		2,
		"",
		"Test rule",
		-1,
		0,
	}
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, testEquipment)
	roster.Units[3].EquipmentsInfo = append(roster.Units[3].EquipmentsInfo, testEquipment)

	err := roster.IsRosterValid()
	if err == ErrorCountEquipmernOnRoster {
		t.Fatal(err)
	}
}

func TestInvalidSoldarRole(t *testing.T) {
	roster := createValidRosterDetailedInfo()
	medic := EquipmentDetailedInfo {
		"0abf8d9c-8c8c-4706-9c97-c85121e3545f",
		"Medic",
		1,
		1,
		"Medic",
		"This soldier is the fireteam is medic. He does not roll to perform medical aid on other soldiers (he must still use 1 order). When this soldier performs medical aid to another soldier that has suffered 3 wounds, the -1 modafire for receiving 3 wounds is removed. This soldier cannot benefit from Medic on himself.",
		-1,
		0,
	}
	roster.Units[1].EquipmentsInfo = append(roster.Units[1].EquipmentsInfo, medic)


	err := roster.IsRosterValid()
	if err != ErrorUnitRole {
		t.Fatal(err)
	}
}