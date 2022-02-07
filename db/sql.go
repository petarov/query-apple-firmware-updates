package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Device struct {
	Id      int
	Product string
	Name    string
}

type DeviceUpdate struct {
	Id         int
	DeviceId   int
	Device     *Device
	BuildId    string
	Version    string
	ReleasedOn string
	Attributes string
}

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

	count, err := importDevices(index)
	if err != nil {
		return fmt.Errorf("Failed adding devices index to database : %w", err)
	}

	fmt.Printf("Added %d devices to database.", count)

	return nil
}

func createSchema() (err error) {
	stmtDevice := `
	CREATE TABLE IF NOT EXISTS Device (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product TEXT UNIQUE NOT NULL,
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
		FOREIGN KEY(device_id) 
		REFERENCES Device(id)
		ON DELETE CASCADE
	);`

	if _, err = db.Exec(stmtUpdate); err != nil {
		return fmt.Errorf("Failed create Update schema : %w", err)
	}

	return nil
}

func importDevices(index *DevicesIndex) (int, error) {
	inserted := 0

	for k, v := range index.revIndex {
		result, err := db.Exec(`INSERT INTO Device(product, name) VALUES($1, $2)`, k, v)
		if err != nil {
			return 0, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		inserted += int(rowsAffected)
	}

	return inserted, nil
}
