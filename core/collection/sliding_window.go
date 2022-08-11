package collection

import (
	"sync"
	"time"
)

var initTime = time.Now().AddDate(-1, 0, 0)

// SlidingWindowOption customize the SlidingWindow.
type SlidingWindowOption func(*SlidingWindow)

// IgnoreCurrentBucket ignore current bucket.
func IgnoreCurrentBucket() SlidingWindowOption {
	return func(s *SlidingWindow) {
		s.ignoreCurrent = true
	}
}

// SlidingWindow defines a Sliding window to calculate the events in buckets with time interval.
type SlidingWindow struct {
	rw            sync.RWMutex
	ignoreCurrent bool
	interval      time.Duration
	lastTime      time.Duration // start time of the last bucket

	offset  int
	size    int
	buckets []*Bucket
}

// NewSlidingWindow returns a SlidingWindow that with size buckets and time interval,
// use opts to customize the SlidingWindow.
func NewSlidingWindow(size int, interval time.Duration, opts ...SlidingWindowOption) *SlidingWindow {
	if size < 1 {
		panic("collection: size must be greater than 0")
	}
	buckets := make([]*Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = new(Bucket)
	}
	w := &SlidingWindow{
		ignoreCurrent: false,
		interval:      interval,
		lastTime:      time.Since(initTime),
		offset:        0,
		size:          size,
		buckets:       buckets,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Add adds value to current bucket.
func (s *SlidingWindow) Add(v float64) *SlidingWindow {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.updateOffset()
	s.buckets[s.offset%s.size].add(v)
	return s
}

// Reduce runs fn on all buckets, ignore current bucket if ignoreCurrent was set.
func (s *SlidingWindow) Reduce(fn func(b *Bucket)) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	var diff int
	span := s.span()
	// ignore current bucket, because of partial data
	if span == 0 && s.ignoreCurrent {
		diff = s.size - 1
	} else {
		diff = s.size - span
	}
	if diff > 0 {
		offset := (s.offset + span + 1) % s.size
		for i := 0; i < diff; i++ {
			fn(s.buckets[(offset+i)%s.size])
		}
	}
}

func (s *SlidingWindow) span() int {
	offset := int((time.Since(initTime) - s.lastTime) / s.interval)
	if offset >= 0 && offset < s.size {
		return offset
	}
	return s.size
}

func (s *SlidingWindow) updateOffset() {
	span := s.span()
	offset := s.offset
	// reset expired buckets
	for i := 0; i < span; i++ {
		s.buckets[(offset+i+1)%s.size].reset()
	}

	s.offset = (offset + span) % s.size
	cur := time.Since(initTime)
	// align to interval time boundary
	s.lastTime = cur - (cur-s.lastTime)%s.interval
}

// Bucket defines the bucket that holds sum and num of additions.
type Bucket struct {
	sum   float64
	count int64
}

func (b *Bucket) add(v float64) {
	b.sum += v
	b.count++
}

func (b *Bucket) reset() {
	b.sum = 0
	b.count = 0
}
func (b *Bucket) Sum() float64 { return b.sum }
func (b *Bucket) Count() int64 { return b.count }
