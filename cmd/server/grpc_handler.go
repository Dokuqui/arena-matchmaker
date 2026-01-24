package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/queue"
	"context"
	"fmt"
)

type GrpcServer struct {
	matchmaker.UnimplementedMatchmakerServiceServer
	QueueClient *queue.Client
}

func (s *GrpcServer) FindMatch(ctx context.Context, req *matchmaker.FindMatchRequest) (*matchmaker.FindMatchResponse, error) {
	fmt.Printf("📥 Received Match Request | Player: %s | MMR: %d\n", req.PlayerId, req.Mmr)

	err := s.QueueClient.AddToQueue(ctx, req.PlayerId, req.Mmr)
	if err != nil {
		return nil, fmt.Errorf("internal redis error: %v", err)
	}

	fakeTicketID := fmt.Sprintf("ticket_%s", req.PlayerId)
	return &matchmaker.FindMatchResponse{
		TicketId: fakeTicketID,
	}, nil
}
