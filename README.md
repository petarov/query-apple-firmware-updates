# Query Apple Devices OS Updates

A web service that delivers firmware update infos about Apple devices.

This is in fact a caching proxy for [ipsw.me](https://ipsw.me/). Update infos are being stored in a local SQLite database.

# API

List of available API calls:

- `/api` - shows all available junctions
- `/api/devices` - Fetches a list of all Apple devices
- `/api/devices/:product` - Fetches a single Apple device by product name
- `/api/updates/:product` - Fetches device updates by product name
    
# Installation

    // TODO

# License 

[MIT](LICENSE)
