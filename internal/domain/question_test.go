package domain

import "testing"

func TestQuestionStats(t *testing.T) {
	stats := QuestionStats{YesAnswers: 8, SometimesAnswers: 1, NeverAnswers: 1, ConfusingVotes: 11}
	if !stats.ShouldArchiveForConfusing(5) {
		t.Fatal("expected archive decision to be true")
	}
	if stats.PercentNormal() != 80 {
		t.Fatalf("expected 80, got %v", stats.PercentNormal())
	}
}
