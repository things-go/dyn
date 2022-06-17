package collection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const duration = time.Millisecond * 50

func TestNewSlidingWindow(t *testing.T) {
	assert.NotNil(t, NewSlidingWindow(10, time.Second))
	assert.Panics(t, func() {
		NewSlidingWindow(0, time.Second)
	})
}

func TestSlidingWindowAdd(t *testing.T) {
	r := NewSlidingWindow(3, duration)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum())
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	elapse()
	r.Add(2).Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	elapse()
	r.Add(4).Add(5).Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
	elapse()
	r.Add(7)
	assert.Equal(t, []float64{5, 15, 7}, listBuckets())
}

func TestSlidingWindowReset(t *testing.T) {
	const size = 3
	r := NewSlidingWindow(size, duration, IgnoreCurrentBucket())
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum())
		})
		return buckets
	}
	r.Add(1)
	elapse()
	assert.Equal(t, []float64{0, 1}, listBuckets())
	elapse()
	assert.Equal(t, []float64{1}, listBuckets())
	elapse()
	assert.Nil(t, listBuckets())

	// cross window
	r.Add(1)
	time.Sleep(duration * 10)
	assert.Nil(t, listBuckets())
}

func TestSlidingWindowReduce(t *testing.T) {
	const size = 4
	tests := []struct {
		win    *SlidingWindow
		expect float64
	}{
		{
			win:    NewSlidingWindow(size, duration),
			expect: 10,
		},
		{
			win:    NewSlidingWindow(size, duration, IgnoreCurrentBucket()),
			expect: 4,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			r := test.win
			for x := 0; x < size; x++ {
				for i := 0; i <= x; i++ {
					r.Add(float64(i))
				}
				if x < size-1 {
					elapse()
				}
			}
			var result float64
			r.Reduce(func(b *Bucket) {
				result += b.Sum()
			})
			assert.Equal(t, test.expect, result)
		})
	}
}

func TestSlidingWindowBucketTimeBoundary(t *testing.T) {
	const size = 3
	interval := time.Millisecond * 30
	r := NewSlidingWindow(size, interval)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum())
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	time.Sleep(time.Millisecond * 45)
	r.Add(2).Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	// sleep time should be less than interval, and make the bucket change happen
	time.Sleep(time.Millisecond * 20)
	r.Add(4).Add(5).Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
	time.Sleep(time.Millisecond * 100)
	r.Add(7).Add(8).Add(9)
	assert.Equal(t, []float64{0, 0, 24}, listBuckets())
}

func TestSlidingWindowDataRace(t *testing.T) {
	const size = 3
	r := NewSlidingWindow(size, duration)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				r.Add(float64(rand.Int63()))
				time.Sleep(duration / 2)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				r.Reduce(func(b *Bucket) {})
			}
		}
	}()
	time.Sleep(duration * 5)
	close(stop)
}

func elapse() {
	time.Sleep(duration + time.Millisecond*2)
}
