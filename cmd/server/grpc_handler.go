package main

import (
	"arena-matchmaker/gen/matchmaker"
	"context"
	"fmt"
	"math/rand"
)

type GrpcServer struct {
	matchmaker.UnimplementedMatchmakerServiceServer
}

func (s *GrpcServer) FindMatch(ctx context.Context, req *matchmaker.FindMatchRequest) (*matchmaker.FindMatchResponse, error) {
	fmt.Printf("📥 Received Match Request | Player: %s | MMR: %d | Region: %s\n",
		req.PlayerId, req.Mmr, req.Region)

	fakeTicketID := fmt.Sprintf("ticket_%d", rand.Intn(10000))

	return &matchmaker.FindMatchResponse{
		TicketId: fakeTicketID,
	}, nil
}
