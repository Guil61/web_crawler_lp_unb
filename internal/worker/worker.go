package worker

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"crawler/internal/models"
	"crawler/internal/parser"
)

type Limiter struct {
	ticker *time.Ticker
}

func NewLimiter(perSecond float64) *Limiter {
	if perSecond <= 0 {
		perSecond = 1
	}
	interval := time.Duration(float64(time.Second) / perSecond)
	return &Limiter{ticker: time.NewTicker(interval)}
}

func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Stop() { l.ticker.Stop() }

type Pool struct {
	workers int
	client  *http.Client
	limiter *Limiter
	stats   *models.Stats

	jobs    <-chan models.Job
	results chan<- models.PageResult

	wg sync.WaitGroup
}

func NewPool(
	workers int,
	jobs <-chan models.Job,
	results chan<- models.PageResult,
	limiter *Limiter,
	stats *models.Stats,
) *Pool {
	return &Pool{
		workers: workers,
		client:  &http.Client{Timeout: 10 * time.Second},
		limiter: limiter,
		stats:   stats,
		jobs:    jobs,
		results: results,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 1; i <= p.workers; i++ {
		p.wg.Add(1)
		go p.run(ctx, i)
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) run(ctx context.Context, id int) {
	defer p.wg.Done()

	for job := range p.jobs {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result := p.process(ctx, id, job)

		select {
		case p.results <- result:
		case <-ctx.Done():
			return
		}
	}
}

func (p *Pool) process(ctx context.Context, id int, job models.Job) models.PageResult {
	res := models.PageResult{Job: job, WorkerID: id}

	if err := p.limiter.Wait(ctx); err != nil {
		res.Err = err
		return res
	}

	fmt.Printf("[Worker %d] Crawling %s (depth %d)\n", id, job.URL, job.Depth)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return p.fail(res, err)
	}
	req.Header.Set("User-Agent", "GoConcurrentCrawler/1.0 (projeto academico)")

	resp, err := p.client.Do(req)
	if err != nil {
		return p.fail(res, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return p.fail(res, fmt.Errorf("status HTTP %d", resp.StatusCode))
	}

	base, err := url.Parse(job.URL)
	if err != nil {
		return p.fail(res, err)
	}

	links, err := parser.ExtractLinks(base, resp.Body)
	if err != nil {
		return p.fail(res, err)
	}

	res.Links = links
	p.stats.AddPagesFetched(1)
	return res
}

func (p *Pool) fail(res models.PageResult, err error) models.PageResult {
	res.Err = err
	p.stats.AddError(1)
	return res
}
