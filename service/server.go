package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/config"
	"github.com/petarov/query-apple-osupdates/db"
)

type ServerContext struct {
	router *http.ServeMux
}

func ServeNow() (err error) {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

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
		config.ListenAddress, config.ListenPort), ctx.router); err != nil {
		return err
	}

	return nil
}
