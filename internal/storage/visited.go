package storage

import "sync"

type VisitedSet struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewVisitedSet() *VisitedSet {
	return &VisitedSet{seen: make(map[string]struct{})}
}

func (v *VisitedSet) Visit(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, ok := v.seen[url]; ok {
		return false
	}
	v.seen[url] = struct{}{}
	return true
}

func (v *VisitedSet) Len() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return len(v.seen)
}
