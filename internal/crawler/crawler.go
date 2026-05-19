package crawler

import (
	"context"
	"fmt"
	"time"

	"crawler/internal/models"
	"crawler/internal/scheduler"
	"crawler/internal/storage"
	"crawler/internal/worker"
)

type Config struct {
	SeedURL   string
	MaxDepth  int
	Workers   int
	RateLimit float64
}

type Crawler struct {
	cfg   Config
	stats *models.Stats
}

func New(cfg Config) *Crawler {
	return &Crawler{cfg: cfg, stats: &models.Stats{}}
}

func (c *Crawler) Run(ctx context.Context) {
	start := time.Now()

	jobs := make(chan models.Job, c.cfg.Workers)
	results := make(chan models.PageResult, c.cfg.Workers)

	visited := storage.NewVisitedSet()
	limiter := worker.NewLimiter(c.cfg.RateLimit)
	defer limiter.Stop()

	pool := worker.NewPool(c.cfg.Workers, jobs, results, limiter, c.stats)
	pool.Start(ctx)

	sch := scheduler.New(c.cfg.MaxDepth, visited, c.stats, jobs, results)

	seed := models.Job{URL: c.cfg.SeedURL, Depth: 0}
	done := make(chan struct{})
	go func() {
		sch.Run(ctx, seed)
		close(done)
	}()

	<-done
	pool.Wait()
	close(results)

	c.printStats(time.Since(start))
}

func (c *Crawler) printStats(elapsed time.Duration) {
	visited, fetched, links, errs := c.stats.Snapshot()

	fmt.Println()
	fmt.Println("======== ESTATÍSTICAS FINAIS ========")
	fmt.Printf("  URLs únicas descobertas : %d\n", visited)
	fmt.Printf("  Páginas baixadas (200)  : %d\n", fetched)
	fmt.Printf("  Links encontrados       : %d\n", links)
	fmt.Printf("  Erros                   : %d\n", errs)
	fmt.Printf("  Tempo total             : %s\n", elapsed.Round(time.Millisecond))
	fmt.Println("=====================================")
}
