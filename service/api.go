package service

import (
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
}

func (api *Api) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := db.FetchAllDevices()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
		} else {
			for _, d := range devices {
				fmt.Fprintf(w, "Device: %d\t%s\t%s\n", d.Id, d.Product, d.Name)
			}
		}
	}
}
