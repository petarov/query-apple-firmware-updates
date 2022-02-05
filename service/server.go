package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/config"
)

type ServerContext struct {
	DevicesIndex map[string]string

	router *http.ServeMux
}

func ServeNow() error {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

	attachApi(ctx)

	fmt.Printf("Serving at %s and port %d ...\n", config.ListenAddress, config.ListenPort)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), ctx.router); err != nil {
		return err
	}

	return nil
}
