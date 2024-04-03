package rate

import (
	"sync"
	"time"
)

// A Limiter controls how frequently events are allowed to happen.
// It implements a "fix window" of size b, initially empty, and then refilled
// Remove at a specific rate
type Limiter struct {
	mu sync.Mutex
	// limit Indicates the maximum number of tokens allowed within the window
	limit      int64
	windowSize time.Duration
	// interval Indicates the rate limit for token generation
	interval time.Duration
	// bucket Used to store the number of tokens that have been used during the window time
	bucket []int64 //TODO: Excessive memory usage
}

func NewLimiter(limit int64, windowSize, interval time.Duration) *Limiter {
	window := windowSize
	if windowSize == 0 {
		window = interval
	}
	return &Limiter{
		limit:      limit,
		windowSize: window,
		interval:   interval,
		bucket:     make([]int64, 0),
	}
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.limit < 0 {
		return true
	}

	current := time.Now().UnixNano()
	lastTime := int64(0)

	l.slide(current)

	if len(l.bucket) > 0 {
		lastTime = l.bucket[len(l.bucket)-1]
	}

	if int64(len(l.bucket)) >= l.limit || current-lastTime < l.interval.Nanoseconds() {
		return false
	}

	l.bucket = append(l.bucket, current)
	return true
}

func (l *Limiter) slide(current int64) {
	target := current - l.windowSize.Nanoseconds()
	index := search(l.bucket, target)
	if index == -1 {
		l.bucket = make([]int64, 0)
	} else {
		l.bucket = l.bucket[index:]
	}
}

func search(arr []int64, target int64) int {
	left, right := 0, len(arr)-1
	if right < 0 || arr[right] < target {
		return -1
	}
	index := -1
	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] >= target {
			index = mid
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return index
}
