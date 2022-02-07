package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/db"
)

const (
	BASE = "/api"
)

type Api struct {
	ctx *ServerContext
}

func attachApi(serverCtx *ServerContext) {
	fmt.Printf("Attaching API junctions at %s ...\n", BASE)

	api := &Api{serverCtx}
	api.ctx.router.HandleFunc("/api/index", api.handleIndex())
	api.ctx.router.HandleFunc("/api/devices", api.handleDevices())

}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := map[string]string{
			"/api":         "This",
			"/api/devices": "Fetch all devices",
		}
		resp, _ := json.Marshal(routes)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (api *Api) handleDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
