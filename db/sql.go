package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/petarov/query-apple-osupdates/client"
)

type Device struct {
	Id      int    `json:"-"`
	Product string `json:"product"`
	Name    string `json:"name"`
}

type DeviceUpdate struct {
	Id         int     `json:"-"`
	DeviceId   int     `json:"-"`
	Device     *Device `json:"device"`
	BuildId    string  `json:"build_id"`
	Version    string  `json:"version"`
	ReleasedOn string  `json:"released_on"`
	Attributes struct {
		IPSW *client.IPSWInfo `json:"ipsw"`
	} `json:"attributes"`
}

var db *sql.DB

func InitDb(path string, jsonDB *DevicesJsonDB) (err error) {
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("Error opening database at %s : %w", path, err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("Failed database ping : %w", err)
	}

	if err = createSchema(); err != nil {
		return err
	}

	count, err := importDevices(jsonDB)
	if err != nil {
		return fmt.Errorf("Failed adding devices index to database : %w", err)
	}

	fmt.Printf("Added %d new device(s) to database.\n", count)

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

func importDevices(jsonDB *DevicesJsonDB) (int, error) {
	inserted := 0

	devices, err := FetchAllDevices()
	if err != nil {
		return 0, err
	}

	lookup := make(map[string]*Device)
	for _, d := range devices {
		lookup[d.Product] = d
	}

	for k, v := range jsonDB.mapping {
		_, ok := lookup[k]
		if !ok {
			_, err := db.Exec(`INSERT INTO Device(product, name) VALUES($1, $2)`, k, v)
			if err != nil {
				return 0, err
			}

			inserted += 1
		}
	}

	return inserted, nil
}

func FetchAllDevices() ([]*Device, error) {
	rows, err := db.Query(`SELECT * FROM Device`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all devices: %w", err)
	}
	defer rows.Close()

	devices := make([]*Device, 0)
	for rows.Next() {
		device := new(Device)
		err := rows.Scan(&device.Id, &device.Product, &device.Name)
		if err != nil {
			return nil, fmt.Errorf("Error scanning device row: %w", err)
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func FetchDeviceByProduct(product string) (*Device, error) {
	device := new(Device)

	err := db.QueryRow(`SELECT * FROM Device WHERE product = $1`, product).Scan(
		&device.Id, &device.Product, &device.Name)
	if err != nil {
		return nil, fmt.Errorf("Error fetching device by product '%s': %w", product, err)
	}

	return device, nil
}
