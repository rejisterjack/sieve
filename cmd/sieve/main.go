package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/rejisterjack/sieve"
	"github.com/rejisterjack/sieve/flaglog"
	"github.com/rejisterjack/sieve/internal/httpapi"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "serve")
	}

	switch os.Args[1] {
	case "serve":
		serveCmd(os.Args[2:])
	case "logs":
		logsCmd(os.Args[2:])
	case "-h", "-help", "--help", "help":
		fmt.Fprintf(os.Stderr, "usage: sieve <serve|logs> [flags]\n")
		os.Exit(0)
	default:
		// Backward compatible: flags without a subcommand imply serve.
		if len(os.Args) > 1 && os.Args[1] != "" && os.Args[1][0] == '-' {
			serveCmd(os.Args[1:])
			return
		}
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	configPath := fs.String("config", "config.example.yaml", "path to YAML config")
	_ = fs.Parse(args)

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

func logsCmd(args []string) {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)
	path := fs.String("path", "sieve.log", "path to JSONL flag log")
	risk := fs.String("risk", "", "filter by risk (suspicious|blocked)")
	category := fs.String("category", "", "filter by match category")
	limit := fs.Int("limit", 0, "max entries to print (0 = all)")
	summary := fs.Bool("summary", false, "print counts by risk and category")
	_ = fs.Parse(args)

	entries, err := flaglog.ReadEntries(*path)
	if err != nil {
		log.Fatalf("read log: %v", err)
	}
	filtered := flaglog.ApplyFilter(entries, flaglog.Filter{
		Risk:     sieve.Risk(*risk),
		Category: *category,
		Limit:    *limit,
	})

	if *summary {
		sum := flaglog.Summarize(filtered)
		fmt.Printf("total: %d\n", sum.Total)
		fmt.Println("by_risk:")
		risks := make([]string, 0, len(sum.ByRisk))
		for r := range sum.ByRisk {
			risks = append(risks, string(r))
		}
		sort.Strings(risks)
		for _, r := range risks {
			fmt.Printf("  %s: %d\n", r, sum.ByRisk[sieve.Risk(r)])
		}
		fmt.Println("by_category:")
		cats := make([]string, 0, len(sum.ByCategory))
		for c := range sum.ByCategory {
			cats = append(cats, c)
		}
		sort.Strings(cats)
		for _, c := range cats {
			fmt.Printf("  %s: %d\n", c, sum.ByCategory[c])
		}
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	for _, e := range filtered {
		if err := enc.Encode(e); err != nil {
			log.Fatal(err)
		}
	}
}
