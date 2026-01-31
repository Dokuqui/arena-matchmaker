package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/logic"
	"arena-matchmaker/pkg/queue"
	"context"
	"fmt"
)

type GrpcServer struct {
	matchmaker.UnimplementedMatchmakerServiceServer
	QueueClient *queue.Client
}

func (s *GrpcServer) FindMatch(ctx context.Context, req *matchmaker.FindMatchRequest) (*matchmaker.FindMatchResponse, error) {
	realMMR, err := s.QueueClient.GetPlayerMMR(ctx, req.PlayerId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player profile: %v", err)
	}

	fmt.Printf("Received Match Request | Player: %s | MMR: %d\n", req.PlayerId, realMMR)

	err = s.QueueClient.AddToQueue(ctx, req.PlayerId, int32(realMMR))
	if err != nil {
		return nil, fmt.Errorf("internal redis error: %v", err)
	}

	fakeTicketID := fmt.Sprintf("ticket_%s", req.PlayerId)
	return &matchmaker.FindMatchResponse{
		TicketId: fakeTicketID,
	}, nil
}

func (s *GrpcServer) ReportResult(ctx context.Context, req *matchmaker.ReportResultRequest) (*matchmaker.ReportResultResponse, error) {
	fmt.Printf("Reporting Result | Match: %s | Winner: %s | Loser: %s\n",
		req.MatchId, req.WinnerId, req.LoserId)

	winnerMMR, err := s.QueueClient.GetPlayerMMR(ctx, req.WinnerId)
	if err != nil {
		return nil, err
	}
	loserMMR, err := s.QueueClient.GetPlayerMMR(ctx, req.LoserId)
	if err != nil {
		return nil, err
	}

	newWinnerMMR, newLoserMMR := logic.CalculateElo(int32(winnerMMR), int32(loserMMR))

	err = s.QueueClient.SetPlayerMMR(ctx, req.WinnerId, int(newWinnerMMR))
	if err != nil {
		return nil, err
	}
	err = s.QueueClient.SetPlayerMMR(ctx, req.LoserId, int(newLoserMMR))
	if err != nil {
		return nil, err
	}

	fmt.Printf("   Start: [%d vs %d] -> End: [%d vs %d]\n",
		winnerMMR, loserMMR, newWinnerMMR, newLoserMMR)

	return &matchmaker.ReportResultResponse{
		WinnerNewMmr: newWinnerMMR,
		LoserNewMmr:  newLoserMMR,
	}, nil
}
