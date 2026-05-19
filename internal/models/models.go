package models

import "sync/atomic"

type Job struct {
	URL   string
	Depth int
}

type PageResult struct {
	Job      Job
	Links    []string
	WorkerID int
	Err      error
}

type Stats struct {
	PagesVisited int64
	PagesFetched int64
	LinksFound   int64
	Errors       int64
}

func (s *Stats) AddPagesVisited(n int64) { atomic.AddInt64(&s.PagesVisited, n) }
func (s *Stats) AddPagesFetched(n int64) { atomic.AddInt64(&s.PagesFetched, n) }
func (s *Stats) AddLinksFound(n int64)   { atomic.AddInt64(&s.LinksFound, n) }
func (s *Stats) AddError(n int64)        { atomic.AddInt64(&s.Errors, n) }

func (s *Stats) Snapshot() (visited, fetched, links, errors int64) {
	return atomic.LoadInt64(&s.PagesVisited),
		atomic.LoadInt64(&s.PagesFetched),
		atomic.LoadInt64(&s.LinksFound),
		atomic.LoadInt64(&s.Errors)
}
