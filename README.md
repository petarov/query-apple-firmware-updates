# Query Apple Devices Firmware Updates

[![CI Build](https://github.com/petarov/query-apple-firmware-updates/actions/workflows/build.yml/badge.svg)](https://github.com/petarov/query-apple-firmware-updates/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/petarov/query-apple-firmware-updates)](https://goreportcard.com/report/github.com/petarov/query-apple-firmware-updates)

A web service that delivers firmware update infos about Apple devices.

This is in fact a caching proxy for [ipsw.me](https://ipsw.me/). Update infos are being stored in a local SQLite database.

# Installation

Download `devices.json` from [SeparateRecords/apple_device_identifiers](https://github.com/SeparateRecords/apple_device_identifiers).

Run `make` to produce binaries in the `dist` folder.

# Usage

To start the service on `[::1]:7095` run:

    ./qadfu_linux_amd64 -devices devices.json -db database.db

The webapp is available at `http://localhost:7095`.

## Run as systemd daemon

Download the binary and `devices.json` file to `/opt/qadfu/` on your server. 

Create a bash `/opt/qadfu/start.sh` startup script:

```bash 
#!/bin/sh

./qadfu_linux_amd64 -addr localhost -devices devices.json -db database.db
```

Create the following file under `/lib/systemd/system/qadfu.service`:

```bash
[Unit]
Description=Query Apple Devices Firmware Updates
After=nginx.service

[Service]
Type=simple
Restart=always
RestartSec=5s
Restart=on-failure
WorkingDirectory=/opt/qadfu
ExecStart=/opt/qadfu/start.sh

[Install]
WantedBy=multi-user.target
```

To enable the systemctl service run:

    systemctl reenable /lib/systemd/system/qadfu.service

To start the service run:

    systemctl start qadfu.service


# API

List of available API calls:

- `/api` - shows all available junctions
- `/api/devices` - Fetches a list of all Apple devices
- `/api/devices/:product` - Fetches a single Apple device by product name
- `/api/devices/search?key=:key` - Fetches a list of matching devices given the `key` parameter
- `/api/updates/:product` - Fetches device updates by product name

# Development

To install deps run:

    go get

To run the server:

    go run -tags "fts5"  main.go -devices devices.json -db database.db

The build tag `fts5` enables the SQLite FTS5 extension in the `mattn/go-sqlite3` lib.

# License 

[MIT](LICENSE)
