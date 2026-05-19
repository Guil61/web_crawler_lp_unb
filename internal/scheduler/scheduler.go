package scheduler

import (
	"context"
	"fmt"

	"crawler/internal/models"
	"crawler/internal/storage"
)

type Scheduler struct {
	maxDepth int
	visited  *storage.VisitedSet
	stats    *models.Stats

	jobs    chan models.Job
	results chan models.PageResult
}

func New(
	maxDepth int,
	visited *storage.VisitedSet,
	stats *models.Stats,
	jobs chan models.Job,
	results chan models.PageResult,
) *Scheduler {
	return &Scheduler{
		maxDepth: maxDepth,
		visited:  visited,
		stats:    stats,
		jobs:     jobs,
		results:  results,
	}
}

func (s *Scheduler) Run(ctx context.Context, seed models.Job) {
	defer close(s.jobs)

	pending := make([]models.Job, 0)

	if s.visited.Visit(seed.URL) {
		pending = append(pending, seed)
		s.stats.AddPagesVisited(1)
	}

	inFlight := 0

	for len(pending) > 0 || inFlight > 0 {
		var sendCh chan<- models.Job
		var next models.Job
		if len(pending) > 0 {
			sendCh = s.jobs
			next = pending[0]
		}

		select {
		case <-ctx.Done():
			fmt.Println("\n[Scheduler] cancelamento recebido — encerrando...")
			return

		case sendCh <- next:
			pending = pending[1:]
			inFlight++

		case res := <-s.results:
			inFlight--
			s.handleResult(res, &pending)
		}
	}
}

func (s *Scheduler) handleResult(res models.PageResult, pending *[]models.Job) {
	if res.Err != nil {
		fmt.Printf("[Scheduler] erro em %s: %v\n", res.Job.URL, res.Err)
		return
	}

	s.stats.AddLinksFound(int64(len(res.Links)))

	nextDepth := res.Job.Depth + 1
	if nextDepth > s.maxDepth {
		return
	}

	novos := 0
	for _, link := range res.Links {
		if s.visited.Visit(link) {
			*pending = append(*pending, models.Job{URL: link, Depth: nextDepth})
			s.stats.AddPagesVisited(1)
			novos++
		}
	}

	if novos > 0 {
		fmt.Printf("[Scheduler] %s -> %d links (%d novos enfileirados)\n",
			res.Job.URL, len(res.Links), novos)
	}
}
