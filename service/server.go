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
	v, _ := ctx.Devices.Get("iPod3,1")
	v2, _ := ctx.Devices.Get("iPod touch (3rd generation)")
	fmt.Printf("DEV: %v\n", v)
	fmt.Printf("DEV: %v\n", v2)

	fmt.Printf("Serving at %s and port %d ...\n", config.ListenAddress, config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), ctx.router); err != nil {
		return err
	}

	return nil
}
