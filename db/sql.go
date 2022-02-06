package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDb(path string, index *DevicesIndex) (err error) {
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("Error opening database at %s : %w", path, err)
	}
	// defer db.Close()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("Failed database ping : %w", err)
	}

	if err = createSchema(); err != nil {
		return err
	}

	importDevices(index)

	return nil
}

func createSchema() (err error) {
	stmtDevice := `
	CREATE TABLE IF NOT EXISTS Device (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product TEXT NOT NULL,
		name TEXT NOT NULL
	);`

	if _, err = db.Exec(stmtDevice); err != nil {
		return fmt.Errorf("Failed create Device schema : %w", err)
	}

	stmtUpdate := `
	CREATE TABLE IF NOT EXISTS DeviceUpdate (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id INTEGER NOT NULL,
		build_id TEXT NOT NULL,
		version  TEXT NOT NULL,
		released_on TEXT NOT NULL,
		attributes TEXT NOT NULL,
		FOREIGN KEY(device_id) REFERENCES Device(id)
	);`

	if _, err = db.Exec(stmtUpdate); err != nil {
		return fmt.Errorf("Failed create Update schema : %w", err)
	}

	return nil
}

func importDevices(index *DevicesIndex) (err error) {
	// TODO
	return nil
}
