package common

import "testing"

func TestCheckAllocationStatus(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{status: "running", want: true},
		{status: "pending", want: true},
		{status: "complete", want: false},
		{status: "failed", want: false},
		{status: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := checkAllocationStatus(tt.status); got != tt.want {
				t.Fatalf("checkAllocationStatus(%q) = %t, want %t", tt.status, got, tt.want)
			}
		})
	}
}
