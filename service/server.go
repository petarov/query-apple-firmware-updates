package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/client"
	"github.com/petarov/query-apple-osupdates/config"
	"github.com/petarov/query-apple-osupdates/db"
)

type ServerContext struct {
	router     *http.ServeMux
	ipswClient *http.Client
}

func ServeNow() (err error) {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

	ctx.ipswClient, err = client.NewIPSWClient()
	if err != nil {
		return err
	}

	jsonDB, err := db.LoadDevices(config.DevicePath)
	if err != nil {
		return err
	}

	if err = db.InitDb(config.DbPath, jsonDB); err != nil {
		return err
	}

	attachApi(ctx)

	fmt.Printf("Serving at %s and port %d ...\n", config.ListenAddress, config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), EnableGzip(ctx.router)); err != nil {
		return err
	}

	return nil
}
