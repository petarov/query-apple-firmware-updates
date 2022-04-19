package service

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func attachWww(serverCtx *ServerContext) {
	fmt.Printf("Attaching web app ...\n")
	html, _ := fs.Sub(staticFiles, "static")
	serverCtx.router.Handle("/", http.FileServer(http.FS(html)))
}
