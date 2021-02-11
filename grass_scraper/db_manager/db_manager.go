package db_manager

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"fmt"
	"sync"
	"time"
)

type DBManager struct {
	User string
	Password string
	Uri string
	Database string
	Initialized bool
	Lock sync.Mutex
}

func (manager *DBManager) Initialize(user string, psw string, uri string, db string) {
	if user == "" {
		panic("user cannot be empty")
	}

	if db == "" {
		panic("db cannot be empty")
	}

	if uri == "" {
		panic("uri cannot be empty")
	}

	manager.User = user
	manager.Database = db
	manager.Password = psw
	manager.Uri = uri
	manager.Initialized = true
}

func (manager *DBManager) connect() *sql.DB {
	if !manager.Initialized {
		panic("DBManager is not yet initialized.")
	}

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", manager.User, manager.Password, manager.Uri, manager.Database))

	if err != nil {
		panic(err)
	}

	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	return conn
}

func (manager *DBManager) TableExists (table string) (bool) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	conn := manager.connect()
	defer conn.Close()

	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1;", table)
	// prepare a timeout to deal with network errors
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return false
	}

	_, err2 := stmt.ExecContext(ctx)

	return err2 == nil
}

func (manager *DBManager) InsertBambooRow(row *BambooEntry, table string) error {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	conn := manager.connect()
	defer conn.Close()

	query := fmt.Sprintf(`INSERT INTO %s (genus_species, num_introductions, is_invasive, disputed_native_range) VALUES (?,?,?,?)`, table)

	// prepare a timeout to deal with network errors
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, row.GenusSpecies, row.NumIntroductions, row.IsInvasive, row.DisputedNativeRange)
	if err != nil {
		log.Printf("Error %s when inserting row into grass_table for species %s", err, row.GenusSpecies)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d bamboo entry inserted for species %s", rows, row.GenusSpecies)
	return nil
}

func (manager *DBManager) InsertGrassRow(row *GrassEntry) error {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	conn := manager.connect()
	defer conn.Close()


	query := `INSERT INTO grass_table(genus_species, is_perennial, is_annual, culm_density, rooting_charactersitic, culm_growth, 
			culm_length_min_cm, culm_length_max_cm, culm_diameter_min_mm, culm_diameter_max_mm, is_woody,
			culm_internode, location_broad, location_narrow, notes) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`


	// prepare a timeout to deal with network errors
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, row.GrassSpecies, row.IsPerennial, row.IsAnnual, row.CulmDensity, row.RootingCharactersitic, row.CulmGrowth,
		row.CulmLengthMinCm, row.CulmLengthMaxCm, row.CulmDiameterMinMm, row.CulmDiameterMaxMm, row.IsWoody, row.CulmInternode,
		row.LocationBroad, row.LocationNarrow, row.Notes)
	if err != nil {
		log.Printf("Error %s when inserting row into grass_table for species %s", err, row.GrassSpecies)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d grass entry inserted for species %s", rows, row.GrassSpecies)
	return nil
}