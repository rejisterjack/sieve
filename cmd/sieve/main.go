package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rejisterjack/sieve"
	"github.com/rejisterjack/sieve/internal/flaglog"
	"github.com/rejisterjack/sieve/internal/httpapi"
)

func main() {
	configPath := flag.String("config", "config.example.yaml", "path to YAML config")
	flag.Parse()

	fileCfg, err := sieve.LoadFile(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	engine, err := sieve.New(fileCfg.ToConfig())
	if err != nil {
		log.Fatalf("init sieve: %v", err)
	}

	logger := flaglog.New(fileCfg.LogPath)
	srv := httpapi.New(engine, logger, fileCfg.WebhookURL)

	fmt.Fprintf(os.Stderr, "sieve listening on %s\n", fileCfg.ListenAddr)
	if err := http.ListenAndServe(fileCfg.ListenAddr, srv.Mux); err != nil {
		log.Fatal(err)
	}
}
