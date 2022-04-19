package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/petarov/query-apple-firmware-updates/client"
	"github.com/petarov/query-apple-firmware-updates/config"
	"github.com/petarov/query-apple-firmware-updates/db"
)

const (
	API_INDEX   = "/api"
	API_DEVICES = API_INDEX + "/devices"
	API_UPDATES = API_INDEX + "/updates"
)

type Api struct {
	ctx *ServerContext
}

func attachApi(serverCtx *ServerContext) {
	fmt.Printf("Attaching API junctions at %s ...\n", API_INDEX)

	api := &Api{serverCtx}
	api.ctx.router.HandleFunc(API_INDEX, api.handleIndex())
	api.ctx.router.HandleFunc(API_INDEX+"/", api.handleIndex())
	api.ctx.router.HandleFunc(API_DEVICES+"/search", api.handleDeviceSearch())
	api.ctx.router.HandleFunc(API_DEVICES+"/search/", api.handleDeviceSearch())
	api.ctx.router.HandleFunc(API_DEVICES, api.handleDevices())
	api.ctx.router.HandleFunc(API_DEVICES+"/", api.handleDevices())
	api.ctx.router.HandleFunc(API_UPDATES+"/", api.handleUpdateInfo())
}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := map[string]string{
			API_INDEX:                        "Shows this",
			API_DEVICES:                      "Fetches a list of all Apple devices",
			API_DEVICES + "/:product":        "Fetches a single Apple device by product name",
			API_DEVICES + "/search?key=:key": "Fetches a list of devices given a key parameter",
			API_UPDATES + "/:product":        "Fetches device updates by product name",
		}
		resp, _ := json.Marshal(routes)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (api *Api) handleDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, API_DEVICES), "/")
		if len(product) > 0 {

			device, err := db.FetchDeviceByProduct(product)
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Error: %v", err)
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %v", err)
			} else {
				resp, _ := json.Marshal(device)
				w.Header().Set("Content-Type", "application/json")
				w.Write(resp)
			}
		} else {

			devices, err := db.FetchAllDevices()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %v", err)
			} else {
				resp, _ := json.Marshal(devices)
				w.Header().Set("Content-Type", "application/json")
				// w.Header().Set("Content-Length", strconv.Itoa(len(resp)))
				w.Write(resp)
			}
		}
	}
}

func (api *Api) handleDeviceSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimSpace(r.URL.Query().Get("key"))
		noResults := false

		if len(key) > 0 {
			devices, err := db.FetchAllDevicesByKey(key)
			if errors.Is(err, sql.ErrNoRows) {
				noResults = true
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %v", err)
			} else {
				resp, _ := json.Marshal(devices)
				w.Header().Set("Content-Type", "application/json")
				w.Write(resp)
			}
		} else {
			noResults = true
		}

		if noResults {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
		}
	}
}

func (api *Api) handleUpdateInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := strings.TrimPrefix(r.URL.Path, API_UPDATES+"/")

		if len(product) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No product name specified!"))
			return
		}

		device, err := db.FetchDeviceUpdatesByProduct(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error fetching updates from database: %v", err)
			return
		}

		if len(device.Updates) == 0 {
			fmt.Printf("Fetching info for product %s ...\n", product)

			ipsw, err := client.IPSWGetInfo(api.ctx.ipswClient, product)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %v", err)
			} else {
				db.AddUpdates(product, ipsw)

				device, err = db.FetchDeviceUpdatesByProduct(product)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error fetching updates from database: %v", err)
					return
				}
			}
		} else {
			// schedule product update, if last check time has expired
			if device.LastCheckedOnParsed.Add(time.Minute * time.Duration(config.DbUpdateRefreshIntervalMins)).Before(time.Now().UTC()) {
				api.ctx.workerPool.QueueJob(&Job{
					[]interface{}{product},

					func(params []interface{}) {
						productParam := params[0].(string)

						fmt.Printf("Updating info for product %s ...\n", productParam)

						ipsw, err := client.IPSWGetInfo(api.ctx.ipswClient, productParam)
						if err == nil {
							db.AddUpdates(productParam, ipsw)
						} else {
							fmt.Println("%w", err)
						}
					},
				})
			}
		}

		resp, _ := json.Marshal(device)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
