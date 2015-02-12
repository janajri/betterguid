package betterguid

import (
	"math/rand"
	"sync"
	"time"
)

const (
	// Modeled after base64 web-safe chars, but ordered by ASCII.
	pushChars = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
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
	var id [8 + 12]byte
	mu.Lock()
	timeMs := time.Now().UnixNano() / 1e6
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
		id[19-i] = pushChars[lastRandChars[i]]
	}
	mu.Unlock()

	// put current time at the beginning
	for i := 7; i >= 0; i-- {
		n := int(timeMs % 64)
		id[i] = pushChars[n]
		timeMs = timeMs / 64
	}
	return string(id[:])
}