package model

type Equipment struct {
	Id string
}

type Unit struct {
	Id string
	Equipments []Equipment
}

type Roster struct {
	Id string
	Name string
	IdUser string
	Units []Unit
}



