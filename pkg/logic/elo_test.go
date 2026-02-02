package logic

import (
	"testing"
)

func TestCalculateElo(t *testing.T) {
	w, l := CalculateElo(1000, 1000)
	if w != 1016 {
		t.Errorf("Expected winner 1016, got %d", w)
	}
	if l != 984 {
		t.Errorf("Expected loser 984, got %d", l)
	}

	w2, _ := CalculateElo(2000, 1000)
	if w2 > 2010 {
		t.Errorf("High rank player gained too much points for easy win: %d", w2)
	}

	w3, _ := CalculateElo(1000, 2000)
	if w3 < 1025 {
		t.Errorf("Low rank player didn't gain enough for upset: %d", w3)
	}
}
