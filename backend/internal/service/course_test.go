package service

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{name: "minutes", t: now.Add(-30 * time.Minute), want: "刚刚"},
		{name: "hours", t: now.Add(-5 * time.Hour), want: "5小时前"},
		{name: "days", t: now.Add(-22 * 24 * time.Hour), want: "22天前"},
		{name: "old", t: now.Add(-35 * 24 * time.Hour), want: now.Add(-35 * 24 * time.Hour).Format("2006-01-02")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relativeTime(tt.t); got != tt.want {
				t.Fatalf("relativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}
