package v1

import (
	"context"
	broadcastingpb "github.com/strideynet/go-play/broadcasting/proto"
	"sync"
)

const bufferSize = 16

type BroadcastManager struct {
	sync.Mutex
	subs map[chan string]struct{}
}

func New() *BroadcastManager {
	return &BroadcastManager{
		Mutex: sync.Mutex{},
		subs:  make(map[chan string]struct{}),
	}
}

func (b *BroadcastManager) Broadcast(text string) {
	b.Lock()
	defer b.Unlock()

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

type Service struct {
	broadcastingpb.UnimplementedBroadcastServer
	bm *BroadcastManager
}

func (s *Service) Send(ctx context.Context, req *broadcastingpb.SendRequest) (*broadcastingpb.SendResponse, error) {
	s.bm.Broadcast(req.Message)

	return &broadcastingpb.SendResponse{}, nil
}

func (s *Service) Subscribe(req *broadcastingpb.SubscribeRequest, ss broadcastingpb.Broadcast_SubscribeServer) error {
	ch, cancel := s.bm.Subscribe()
	defer cancel()

	for {
		select {
		case <-ss.Context().Done():
			break
		case msg := <-ch:
			if err := ss.Send(&broadcastingpb.SubscribeResponse{Message: msg}); err != nil {
				return err
			}
		}
	}
}
