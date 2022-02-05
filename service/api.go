package service

import (
	"fmt"
	"net/http"
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
		fmt.Fprintf(w, "Hello World")
	}
}
