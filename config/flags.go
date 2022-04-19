package config

var (
	// Server listen address
	ListenAddress string
	// Server listen port
	ListenPort int
	// Devices index file path
	DevicePath string
	// SQLite Database file path
	DbPath string
	// How often to refresh updates in the database, in minutes
	DbUpdateRefreshIntervalMins int
)
