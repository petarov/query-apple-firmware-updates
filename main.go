package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/petarov/query-apple-osupdates/config"
	"github.com/petarov/query-apple-osupdates/service"
)

const (
	DEFAULT_PORT = 7095
	HEART        = "\u2764"
)

func init() {
	flag.StringVar(&config.ListenAddress, "addr", "[::1]", "Server listen address")
	flag.IntVar(&config.ListenPort, "port", DEFAULT_PORT, "Server listen port")
	flag.StringVar(&config.DevicePath, "devices", "", "Path to devices index registry JSON file")
	flag.StringVar(&config.DbPath, "db", "", "Path to SQLite database file")
}

func verifyPath(path string, what string, mustExist bool) {
	if len(path) < 1 {
		fmt.Printf("Error: %s path not specified!\n\n", what)
		flag.PrintDefaults()
		os.Exit(1)
	}
	if mustExist {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Error: %s path not found at '%s'!\n\n", what, path)
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
}

func main() {
	fmt.Printf("%s Query Apple OS Updates service v%s %s\n\n", HEART, config.VERSION, HEART)

	flag.Parse()
	verifyPath(config.DevicePath, "Devices JSON file", true)
	verifyPath(config.DbPath, "Database file", false)

	if err := service.ServeNow(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
}
