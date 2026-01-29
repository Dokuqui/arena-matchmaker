package logic

import "math"

const KFactor = 32

func CalculateElo(winnerMMR, loserMMR int32) (int32, int32) {
	// Formula: 1 / (1 + 10^((RatingB - RatingA) / 400))
	expectedWin := 1.0 / (1.0 + math.Pow(10.0, float64(loserMMR-winnerMMR)/400.0))

	// Winner gets points based on how "hard" the match was
	newWinnerMMR := float64(winnerMMR) + float64(KFactor)*(1.0-expectedWin)
	newLoserMMR := float64(loserMMR) + float64(KFactor)*(0.0-(1.0-expectedWin))

	return int32(newWinnerMMR), int32(newLoserMMR)
}
