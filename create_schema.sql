CREATE DATABASE IF NOT EXISTS grass_db;

USE grass_db;

CREATE TABLE IF NOT EXISTS grass_table (
	genus_species varchar(255) PRIMARY KEY,
	is_perennial bool,
	culm_density varchar(255),
	rooting_charactersitic varchar(255),
	culm_growth varchar(255),
	culm_length_min_cm int,
	culm_length_max_cm int,
	culm_diameter_min_mm int,
	culm_diameter_max_mm int,
	is_woody bool,
	culm_internode varchar(255),
    location_broad varchar(255),
    location_narrow varchar(255),
    notes varchar(255)
);

CREATE USER IF NOT EXISTS grass_user IDENTIFIED BY 'the_grass_user';
GRANT ALL PRIVILEGES ON grass_db.* TO grass_user@localhost;