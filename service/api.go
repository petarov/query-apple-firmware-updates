package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/petarov/query-apple-osupdates/db"
)

const (
	API_INDEX   = "/api"
	API_DEVICES = API_INDEX + "/devices"
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
}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := map[string]string{
			API_INDEX:                 "Shows this",
			API_DEVICES:               "Fetches a list of all Apple devices",
			API_DEVICES + "/:product": "Fetches a single Apple device by its product name",
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
				w.Write(resp)
			}
		}
	}
}
