package betterguid

import (
	"math/rand"
	"sync"
	"time"
)

const (
	// Modeled after base64 web-safe chars, but ordered by ASCII.
	pushChars        = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	pushCharsReverse = "zyxwvutsrqponmlkjihgfedcba_ZYXWVUTSRQPONMLKJIHGFEDCBA9876543210-"
)

var (
	// Timestamp of last push, used to prevent local collisions if you push twice in one ms.
	lastPushTimeMs int64
	// We generate 72-bits of randomness which get turned into 12 characters and appended to the
	// timestamp to prevent collisions with other clients.  We store the last characters we
	// generated because in the event of a collision, we'll use those same characters except
	// "incremented" by one.
	lastRandChars [12]int
	mu            sync.Mutex
	rnd           *rand.Rand
)

func init() {
	// have to seed to get randomness
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 12; i++ {
		lastRandChars[i] = rnd.Intn(64)
	}
}

// New creates a new guid.
func New() string {
	timeMs := time.Now().UTC().UnixNano() / 1e6
	return generate(timeMs, true)
}

// Ascending creates a new ascending guid
func Ascending() string {
	timeMs := time.Now().UTC().UnixNano() / 1e6
	return generate(timeMs, true)
}

// AscendingFrom creates a new ascending guid based on a provided time in ms
func AscendingFrom(timeMs int64) string {
	return generate(timeMs, true)
}

// Descending creates a new descending guid
func Descending() string {
	timeMs := time.Now().UTC().UnixNano() / 1e6
	return generate(timeMs, false)
}

// DescendingFrom creates a new descending guid based on a provided time in ms
func DescendingFrom(timeMs int64) string {
	return generate(timeMs, false)
}

func generate(timeMs int64, ascending bool) string {
	pool := pushChars
	if !ascending {
		pool = pushCharsReverse
	}
	var id [8 + 12]byte
	mu.Lock()
	if timeMs == lastPushTimeMs {
		// increment lastRandChars
		for i := 0; i < 12; i++ {
			lastRandChars[i]++
			if lastRandChars[i] < 64 {
				break
			}
			// increment the next byte
			lastRandChars[i] = 0
		}
	}
	lastPushTimeMs = timeMs
	// put random as the second part
	for i := 0; i < 12; i++ {
		id[19-i] = pool[lastRandChars[i]]
	}
	mu.Unlock()
	// put current time at the beginning
	for i := 7; i >= 0; i-- {
		n := int(timeMs % 64)
		id[i] = pool[n]
		timeMs = timeMs / 64
	}
	return string(id[:])
}
