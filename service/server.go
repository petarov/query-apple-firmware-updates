package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/config"
	"github.com/petarov/query-apple-osupdates/db"
)

type ServerContext struct {
	DevicesIndex map[string]string

	router *http.ServeMux
}

func ServeNow() (err error) {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

	attachApi(ctx)

	ctx.DevicesIndex, err = db.LoadDevices(config.DevicePath)
	if err != nil {
		return err
	}
	fmt.Printf("DEV: %v\n", ctx.DevicesIndex)

	fmt.Printf("Serving at %s and port %d ...\n", config.ListenAddress, config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), ctx.router); err != nil {
		return err
	}

	return nil
}
