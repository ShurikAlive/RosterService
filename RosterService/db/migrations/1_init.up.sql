CREATE TABLE `roster_db`.`rosters` (
    id varchar(50),
	Name varchar(150),
	idUser varchar(50),
	status int
);

CREATE TABLE `roster_db`.`roster_units` (
    idRoster varchar(50),
	idUnit varchar(50)
);

CREATE TABLE `roster_db`.`roster_equipments` (
    idRoster varchar(50),
	idUnit varchar(50),
    idEquipment varchar(50)
);

