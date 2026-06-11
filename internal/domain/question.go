package domain

type AnswerChoice string

const (
	AnswerYes       AnswerChoice = "yes_for_me"
	AnswerSometimes AnswerChoice = "sometimes_for_me"
	AnswerNever     AnswerChoice = "never_for_me"
)

type QuestionStats struct {
	YesAnswers       int `json:"yesAnswers"`
	SometimesAnswers int `json:"sometimesAnswers"`
	NeverAnswers     int `json:"neverAnswers"`
	ConfusingVotes   int `json:"confusingVotes"`
}

func (s QuestionStats) TotalAnswers() int {
	return s.YesAnswers + s.SometimesAnswers + s.NeverAnswers
}

func (s QuestionStats) PercentNormal() float64 {
	total := s.TotalAnswers()
	if total == 0 {
		return 0
	}
	return float64(s.YesAnswers) / float64(total) * 100
}

func (s QuestionStats) ShouldArchiveForConfusing(minVotes int) bool {
	if minVotes <= 0 {
		minVotes = 5
	}
	return s.ConfusingVotes >= minVotes && s.ConfusingVotes > s.TotalAnswers()
}
