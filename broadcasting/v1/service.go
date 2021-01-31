package v1

import (
	"context"
	broadcastingpb "github.com/strideynet/go-play/broadcasting/proto"
)

// GRPC Service demonstrates how this BroadcastManager would be used.
type Service struct {
	broadcastingpb.UnimplementedBroadcastServer
	bm *BroadcastManager
}

func (s *Service) Send(ctx context.Context, req *broadcastingpb.SendRequest) (*broadcastingpb.SendResponse, error) {
	s.bm.Broadcast(req.Message)

	return &broadcastingpb.SendResponse{}, nil
}

func (s *Service) Subscribe(req *broadcastingpb.SubscribeRequest, ss broadcastingpb.Broadcast_SubscribeServer) error {
	// The BM interface is easy to understand, it just returns a channel that can then be consumed as part of a for select
	// loop.
	ch, cancel := s.bm.Subscribe()
	defer cancel() // If a ``defer cancel()`` is forgotten by a library user, this will cause a memory leak.

	for {
		select {
		// We need to ensure we watch for cancellation, otherwise this will not close properly if a client disconnects.
		case <-ss.Context().Done():
			break
		case msg := <-ch:
			if err := ss.Send(&broadcastingpb.SubscribeResponse{Message: msg}); err != nil {
				return err
			}
		}
	}
}
