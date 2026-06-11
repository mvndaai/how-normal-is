package domain

import "time"

type CommentAction string

const (
	CommentApprove      CommentAction = "approve"
	CommentReject       CommentAction = "reject"
	CommentStrongReject CommentAction = "strong_reject"
)

type CommentBanWindow struct {
	Until   time.Time `json:"until,omitempty"`
	Forever bool      `json:"forever"`
	Reason  string    `json:"reason"`
}

func BanWindowForStrongRejects(now time.Time, strongRejectCount int) CommentBanWindow {
	switch {
	case strongRejectCount >= 6:
		return CommentBanWindow{Forever: true, Reason: "6+ strong rejects"}
	case strongRejectCount == 5:
		return CommentBanWindow{Until: now.AddDate(1, 0, 0), Reason: "5 strong rejects"}
	case strongRejectCount == 4:
		return CommentBanWindow{Until: now.AddDate(0, 1, 0), Reason: "4 strong rejects"}
	case strongRejectCount == 3:
		return CommentBanWindow{Until: now.AddDate(0, 0, 7), Reason: "3 strong rejects"}
	default:
		return CommentBanWindow{}
	}
}

func ShouldWarnModerator(approved, laterFlagged int) bool {
	if approved < 10 {
		return false
	}
	ratio := float64(laterFlagged) / float64(approved)
	return ratio >= 0.2
}

func ShouldRevokeModerator(approved, laterFlagged int) bool {
	if approved < 10 {
		return false
	}
	ratio := float64(laterFlagged) / float64(approved)
	return ratio >= 0.35
}
