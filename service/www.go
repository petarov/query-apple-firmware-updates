package service

import (
	"fmt"
	"net/http"

	"github.com/petarov/query-apple-osupdates/config"
)

func attachWww(serverCtx *ServerContext) {
	fmt.Printf("Attaching web app ...\n")
	serverCtx.router.Handle("/", http.FileServer(http.Dir(config.WebAppPath)))
}
