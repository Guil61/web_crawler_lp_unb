package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"crawler/internal/crawler"
)

func main() {
	urlFlag := flag.String("url", "", "URL inicial (semente) do crawl [obrigatório]")
	depth := flag.Int("depth", 2, "profundidade máxima de navegação")
	workers := flag.Int("workers", 5, "número de workers concorrentes")
	rate := flag.Float64("rate", 10, "limite de requisições por segundo (global)")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Fprintln(os.Stderr, "erro: a flag -url é obrigatória")
		flag.Usage()
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := crawler.Config{
		SeedURL:   *urlFlag,
		MaxDepth:  *depth,
		Workers:   *workers,
		RateLimit: *rate,
	}

	fmt.Printf("Iniciando crawl: url=%s | depth=%d | workers=%d | rate=%.1f req/s\n\n",
		cfg.SeedURL, cfg.MaxDepth, cfg.Workers, cfg.RateLimit)

	crawler.New(cfg).Run(ctx)
}
