package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/petarov/query-apple-osupdates/client"
	"github.com/petarov/query-apple-osupdates/db"
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
	api.ctx.router.HandleFunc(API_DEVICES, api.handleDevices())
	api.ctx.router.HandleFunc(API_DEVICES+"/", api.handleDevices())
	api.ctx.router.HandleFunc(API_UPDATES+"/", api.handleUpdateInfo())
}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := map[string]string{
			API_INDEX:                 "Shows this",
			API_DEVICES:               "Fetches a list of all Apple devices",
			API_DEVICES + "/:product": "Fetches a single Apple device by product name",
			API_UPDATES + "/:product": "Fetches device updates by product name",
		}
		resp, _ := json.Marshal(routes)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (api *Api) handleDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := strings.TrimPrefix(r.URL.Path, API_DEVICES+"/")
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
			// TODO: schedule update
			api.ctx.workerPool.QueueJob(&Job{
				nil,
				func() {
					fmt.Println("UPDATING ... ")
				},
			})
		}

		resp, _ := json.Marshal(device)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
