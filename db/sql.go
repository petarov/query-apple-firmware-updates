package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/petarov/query-apple-firmware-updates/client"
)

type Device struct {
	Id                  int             `json:"-"`
	Product             string          `json:"product"`
	Name                string          `json:"name"`
	LastCheckedOn       string          `json:"last_checked_on"`
	LastCheckedOnParsed time.Time       `json:"-"`
	Updates             []*DeviceUpdate `json:"updates,omitempty"`
}

type DeviceUpdate struct {
	Id         int              `json:"-"`
	DeviceId   int              `json:"-"`
	BuildId    string           `json:"build_id"`
	Version    string           `json:"version"`
	ReleasedOn string           `json:"released_on"`
	Attributes *client.IPSWInfo `json:"attributes,omitempty"`
}

const DATE_TIME_LAYOUT = "2006-01-02T15:04:05.000Z"

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

	count, err := addDevices(jsonDB)
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
		name TEXT NOT NULL,
		last_checked_on TEXT NOT NULL
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

	if _, err = db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS DeviceFTSI USING fts5(id, name)`); err != nil {
		return fmt.Errorf("Failed create Device full-text index VT : %w", err)
	}

	return nil
}

func createFTSIndex() (err error) {
	if _, err = db.Exec(`INSERT INTO DeviceFTSI SELECT id, name FROM Device`); err != nil {
		return fmt.Errorf("Failed insert devices into full-text index VT : %w", err)
	}

	return nil
}

func addDevices(jsonDB *DevicesJsonDB) (int, error) {
	devices, err := FetchAllDevices()
	if err != nil {
		return 0, err
	}

	lookup := make(map[string]*Device)
	for _, d := range devices {
		lookup[d.Product] = d
	}

	inserted := 0

	now := time.Now().UTC().Format(DATE_TIME_LAYOUT)

	for k, v := range jsonDB.mapping {
		_, ok := lookup[k]
		if !ok {
			res, err := db.Exec(`INSERT INTO Device(product, name, last_checked_on) VALUES($1, $2, $3)`, k, v, now)
			if err != nil {
				return 0, err
			}

			rowId, _ := res.LastInsertId()
			_, err = db.Exec(`INSERT INTO DeviceFTSI(id, name) VALUES($1, $2)`, rowId, v)
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
		err := rows.Scan(&device.Id, &device.Product, &device.Name, &device.LastCheckedOn)
		if err != nil {
			return nil, fmt.Errorf("Error scanning device row: %w", err)
		}

		device.LastCheckedOnParsed, err = time.Parse(DATE_TIME_LAYOUT, device.LastCheckedOn)
		if err != nil {
			return nil, fmt.Errorf("Error parsing last update check date time '%s' : %w", device.LastCheckedOn, err)
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func FetchAllDevicesByKey(key string) ([]*Device, error) {
	rows, err := db.Query(`SELECT * FROM Device 
	WHERE product LIKE $1||'%' OR product LIKE '%'||$2||'%' OR name LIKE '%'||$3||'%'`, key, key, key)
	if err != nil {
		return nil, fmt.Errorf("Error searching devices: %w", err)
	}
	defer rows.Close()

	devices := make([]*Device, 0)
	for rows.Next() {
		device := new(Device)
		err := rows.Scan(&device.Id, &device.Product, &device.Name, &device.LastCheckedOn)
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
		&device.Id, &device.Product, &device.Name, &device.LastCheckedOn)
	if err != nil {
		return nil, fmt.Errorf("Error fetching device by product '%s': %w", product, err)
	}

	device.LastCheckedOnParsed, err = time.Parse(DATE_TIME_LAYOUT, device.LastCheckedOn)
	if err != nil {
		return nil, fmt.Errorf("Error parsing last update check date time '%s' : %w", device.LastCheckedOn, err)
	}

	return device, nil
}

func FetchDeviceUpdatesByProduct(product string) (*Device, error) {
	rows, err := db.Query(`
	SELECT d.*, du.* FROM DeviceUpdate AS du
	LEFT JOIN Device AS d ON du.device_id = d.id
	WHERE d.product = $1 
	ORDER BY DATETIME(du.released_on) DESC`, product)
	if err != nil {
		return nil, fmt.Errorf("Error fetching device updates for %s: %w", product, err)
	}
	defer rows.Close()

	device := new(Device)
	device.Updates = make([]*DeviceUpdate, 0)

	for rows.Next() {
		attributes := ""
		update := new(DeviceUpdate)

		err := rows.Scan(&device.Id, &device.Product, &device.Name, &device.LastCheckedOn,
			&update.Id, &update.DeviceId, &update.BuildId, &update.Version, &update.ReleasedOn, &attributes)
		if err != nil {
			return nil, fmt.Errorf("Error scanning device row: %w", err)
		}

		err = json.Unmarshal([]byte(attributes), &update.Attributes)
		if err != nil {
			return nil, err
		}

		device.Updates = append(device.Updates, update)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(device.Updates) > 0 {
		device.LastCheckedOnParsed, err = time.Parse(DATE_TIME_LAYOUT, device.LastCheckedOn)
		if err != nil {
			return nil, fmt.Errorf("Error parsing last update check date time '%s' : %w", device.LastCheckedOn, err)
		}
	}

	return device, nil
}

func AddUpdates(product string, updatesInfo []*client.IPSWInfo) (int, error) {
	device, err := FetchDeviceByProduct(product)
	if err != nil {
		return 0, err
	}

	inserted := 0

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`INSERT INTO DeviceUpdate(device_id, build_id, version, released_on, attributes) 
	VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}

	for _, v := range updatesInfo {
		attributes, err := json.Marshal(v)
		if err != nil {
			return 0, err
		}

		_, err = stmt.Exec(device.Id, v.BuildId, v.Version, v.ReleaseDate, attributes)
		if err != nil {
			return 0, fmt.Errorf("Error adding device update to database: %w", err)
		}

		inserted += 1
	}

	stmt.Close()

	now := time.Now().UTC().Format(DATE_TIME_LAYOUT)

	_, err = tx.Exec(`UPDATE Device SET last_checked_on=$1 WHERE id=$2`, now, device.Id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return inserted, nil
}
