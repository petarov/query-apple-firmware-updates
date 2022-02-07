package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/config"
	"github.com/petarov/query-apple-osupdates/db"
)

type ServerContext struct {
	Devices *db.DevicesIndex

	router *http.ServeMux
}

func ServeNow() (err error) {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

	attachApi(ctx)

	ctx.Devices, err = db.LoadDevices(config.DevicePath)
	if err != nil {
		return err
	}

	if err = db.InitDb(config.DbPath, ctx.Devices); err != nil {
		return err
	}

	fmt.Printf("Serving at %s and port %d ...\n", config.ListenAddress, config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), ctx.router); err != nil {
		return err
	}

	return nil
}
