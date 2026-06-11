package domain

import (
	"testing"
	"time"
)

func TestBanWindowForStrongRejects(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if bw := BanWindowForStrongRejects(now, 2); bw.Forever || !bw.Until.IsZero() {
		t.Fatalf("unexpected ban window for 2 rejects: %#v", bw)
	}
	if bw := BanWindowForStrongRejects(now, 3); bw.Until.Sub(now) < (6 * 24 * time.Hour) {
		t.Fatalf("expected ~1 week ban, got %#v", bw)
	}
	if bw := BanWindowForStrongRejects(now, 6); !bw.Forever {
		t.Fatalf("expected forever ban, got %#v", bw)
	}
}
