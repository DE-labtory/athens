package main

import (
	"flag"
	"fmt"
	"github.com/DE-labtory/athens"
	"github.com/DE-labtory/cleisthenes"
	"github.com/DE-labtory/cleisthenes/config"
	"github.com/DE-labtory/cleisthenes/core"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
)

func main() {
	host := flag.String("address", "127.0.0.1", "Application address")
	port := flag.Int("port", 18080, "Application port")
	configPath := flag.String("config", "", "User defined config path")

	address := fmt.Sprintf("%s:%d", *host, *port)

	kitLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC)
	httpLogger := kitlog.With(kitLogger, "component", "http")

	config.Init(*configPath)

	txValidator := func(tx cleisthenes.Transaction) bool {
		return true
	}

	node, err := core.New(txValidator)
	if err != nil {
		panic(fmt.Sprintf("Cleisthenes instantiate failed with err: %s", err))
	}

	go func() {
		httpLogger.Log("message", "hbbft started")
		node.Run()
	}()

	httpLogger.Log("message", fmt.Sprintf("http server started: %s", address))
	if err := http.ListenAndServe(address, athens.NewApiHandler(node, httpLogger)); err != nil {
		httpLogger.Log("message", fmt.Sprintf("http server closed: %s", err))
	}
}
