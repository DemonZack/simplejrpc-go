package main

import (
	"path/filepath"

	"simplejrpc-go/core"
	"simplejrpc-go/core/config"
	"simplejrpc-go/example/app/server"
	"simplejrpc-go/os/gpath"
)

func main() {
	env := "test"
	fullPath, _ := filepath.Abs(filepath.Dir("."))
	gpath.GmCfgPath = filepath.Join(filepath.Dir(fullPath), "..")
	core.InitContainer(config.WithConfigEnvFormatterOptionFunc(env))

	appServer := server.NewAppServer()
	appServer.Run()
}
