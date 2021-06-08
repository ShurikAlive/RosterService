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

CREATE TABLE `roster_db`.`roster_equipments` (
    idRoster varchar(50),
	idUnit varchar(50),
    idEquipment varchar(50)
);

CREATE TABLE `roster_db`.`event_execution_tasks` (
	idEvent varchar(50),
	essence varchar(50),
	typeEvent varchar(50),
	idRecord varchar(50),
	creation_time DATETIME DEFAULT CURRENT_TIMESTAMP
);