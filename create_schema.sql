CREATE DATABASE IF NOT EXISTS grass_db;

USE grass_db;

CREATE TABLE IF NOT EXISTS grass_table (
	genus_species varchar(255) PRIMARY KEY,
	is_perennial bool,
    is_annual bool,
	culm_density varchar(1024),
	rooting_charactersitic varchar(1024),
	culm_growth varchar(1024),
	culm_length_min_cm int,
	culm_length_max_cm int,
	culm_diameter_min_mm int,
	culm_diameter_max_mm int,
	is_woody bool,
	culm_internode varchar(1024),
    location_broad varchar(1024),
    location_narrow varchar(1024),
    notes varchar(1024)
);

CREATE USER IF NOT EXISTS grass_user@localhost IDENTIFIED BY 'the_grass_user';
GRANT ALL PRIVILEGES ON grass_db TO grass_user@localhost;
GRANT ALL PRIVILEGES ON grass_db.* TO grass_user@localhost;