# Query Apple Devices OS Updates

A web service that delivers firmware update infos about Apple devices.

This is in fact a caching proxy for [ipsw.me](https://ipsw.me/). Update infos are being stored in a local SQLite database.

# Installation

Download `devices.json` from [SeparateRecords/apple_device_identifiers](https://github.com/SeparateRecords/apple_device_identifiers).

Run `make` to produce binaries in the `dist` folder.

# Usage

To start the service on `[::1]:7095` run:

    ./qaous_linux_amd64 -db database.db -devices devices.json

# API

List of available API calls:

- `/api` - shows all available junctions
- `/api/devices` - Fetches a list of all Apple devices
- `/api/devices/:product` - Fetches a single Apple device by product name
- `/api/updates/:product` - Fetches device updates by product name

# License 

[MIT](LICENSE)
