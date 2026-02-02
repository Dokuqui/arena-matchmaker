package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/logic"
	"arena-matchmaker/pkg/queue"
	"context"
	"fmt"
	"time"
)

type GrpcServer struct {
	matchmaker.UnimplementedMatchmakerServiceServer
	QueueClient *queue.Client
}

func (s *GrpcServer) FindMatch(ctx context.Context, req *matchmaker.FindMatchRequest) (*matchmaker.FindMatchResponse, error) {
	if len(req.PlayerIds) == 0 {
		return nil, fmt.Errorf("player list cannot be empty")
	}
	region := req.Region
	if region == "" {
		region = "EU-WEST"
	}

	ticketID := fmt.Sprintf("ticket_%s_%d", req.PlayerIds[0], time.Now().UnixNano())

	fmt.Printf("📥 Queue Request | Ticket: %s | Region: %s | MMR: %d\n",
		ticketID, region, req.Mmr)

	err := s.QueueClient.AddToQueue(ctx, ticketID, req.PlayerIds, int(req.Mmr), region)
	if err != nil {
		return nil, fmt.Errorf("internal redis error: %v", err)
	}

	return &matchmaker.FindMatchResponse{TicketId: ticketID}, nil
}

func (s *GrpcServer) ReportResult(ctx context.Context, req *matchmaker.ReportResultRequest) (*matchmaker.ReportResultResponse, error) {
	fmt.Printf("📝 Reporting Result | Match: %s | Winners: %v | Losers: %v\n",
		req.MatchId, req.WinnerIds, req.LoserIds)

	for i := range req.WinnerIds {
		if i >= len(req.LoserIds) {
			break
		}

		wID := req.WinnerIds[i]
		lID := req.LoserIds[i]

		wMMR, _ := s.QueueClient.GetPlayerMMR(ctx, wID)
		lMMR, _ := s.QueueClient.GetPlayerMMR(ctx, lID)

		nw, nl := logic.CalculateElo(int32(wMMR), int32(lMMR))

		_ = s.QueueClient.SetPlayerMMR(ctx, wID, int(nw))
		_ = s.QueueClient.SetPlayerMMR(ctx, lID, int(nl))

		fmt.Printf("   Update: %s (%d->%d) vs %s (%d->%d)\n", wID, wMMR, nw, lID, lMMR, nl)
	}

	return &matchmaker.ReportResultResponse{Success: true}, nil
}
