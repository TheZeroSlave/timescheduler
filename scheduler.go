package timescheduler

import (
	"container/heap"
	"context"
	"math"
	"sync"
	"time"
)

type entry struct {
	t time.Time
	f func()
}

type entryQueue []entry

func (pq entryQueue) Len() int { return len(pq) }

func (pq entryQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].t.Before(pq[j].t)
}

func (pq *entryQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func (pq *entryQueue) Push(x any) {
	item := x.(entry)
	*pq = append(*pq, item)
}

func (pq entryQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

type Scheduler struct {
	schedulers entryQueue
	timer      *time.Timer
	mx         *sync.Mutex

	ctx context.Context
}

// Creates a new scheduler
func NewScheduler(ctx context.Context) *Scheduler {
	s := &Scheduler{
		mx:    &sync.Mutex{},
		timer: time.NewTimer(time.Duration(math.MaxInt64)),
		ctx:   ctx,
	}
	go func() {
		for {
			select {
			case <-s.timer.C:
				s.mx.Lock()
				for len(s.schedulers) > 0 {
					recent := s.schedulers[0]
					if time.Until(recent.t) <= 0 {
						recent.f()
						heap.Pop(&s.schedulers)
					} else {
						break
					}
				}

				if len(s.schedulers) != 0 {
					s.timer.Stop()
					s.timer.Reset(time.Until(s.schedulers[0].t))
				}
				s.mx.Unlock()
			case <-s.ctx.Done():
				s.timer.Stop()
				return
			}
		}
	}()
	return s
}

func (s *Scheduler) Add(t time.Time, f func()) {
	if time.Until(t) < 0 {
		return
	}
	s.mx.Lock()
	var first time.Time
	if len(s.schedulers) != 0 {
		first = s.schedulers[0].t
	}
	heap.Push(&s.schedulers, entry{t, f})
	if first != s.schedulers[0].t {
		s.timer.Stop()
		s.timer.Reset(time.Until(t))
	}
	s.mx.Unlock()
}
