package v1

import (
	"sync"
)

// bufferSize dictates the size of the channels returned by Subscribe.
const bufferSize = 16

type BroadcastManager struct {
	// Single mutex to synchronise access to subs map
	sync.RWMutex
	// Whilst this could be a slice, it's significantly cleaner for us to do this as a map since we can easily remove
	// elements by a reference, rather than having to iterate through a slice to find said element.
	subs map[chan string]struct{}
}

// New returns a new BroadcastManager with its fields
func New() *BroadcastManager {
	return &BroadcastManager{
		RWMutex: sync.RWMutex{},
		subs:    make(map[chan string]struct{}),
	}
}

func (b *BroadcastManager) Broadcast(text string) {
	b.RLock()
	defer b.RUnlock()

	// DANGER! If one of these channels is full, and its consumer isn't consuming, this is going to lead to the entire
	// function blocking.
	for ch := range b.subs {
		ch <- text
	}
}

// Subscribe returns a channel that is subscribed to the broadcasts, it also returns a function that removes the
// channel from the subscriptions.
func (b *BroadcastManager) Subscribe() (chan string, func()) {
	b.Lock()
	defer b.Unlock()

	ch := make(chan string, bufferSize)
	b.subs[ch] = struct{}{}

	// Returns function for unregistering, this can be deferred in the caller of Subscribe().
	return ch, func() {
		b.Lock()
		defer b.Unlock()

		delete(b.subs, ch)
	}
}
