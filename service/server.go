package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-firmware-updates/client"
	"github.com/petarov/query-apple-firmware-updates/config"
	"github.com/petarov/query-apple-firmware-updates/db"
)

type ServerContext struct {
	router     *http.ServeMux
	ipswClient *http.Client
	workerPool *WokerPool
}

func ServeNow() (err error) {
	ctx := new(ServerContext)
	ctx.router = http.NewServeMux()

	// ctx.router.HandleFunc("/debug/pprof/", pprof.Index)
	// ctx.router.HandleFunc("/debug/pprof/{action}", pprof.Index)
	// ctx.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

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

	ctx.workerPool = NewWorkPool()
	ctx.workerPool.Start()

	attachApi(ctx)
	attachWww(ctx)

	fmt.Printf("Serving at %s and port %d ...\n", config.ListenAddress, config.ListenPort)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d",
		config.ListenAddress, config.ListenPort), EnableGzip(ctx.router)); err != nil {
		return err
	}

	return nil
}
