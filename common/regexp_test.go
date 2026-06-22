package common

import "testing"

func TestExtractUUID(t *testing.T) {
	const allocID = "7aeb1f6d-1ee7-40e4-b3f0-1fd5248bade8"

	tests := []struct {
		name      string
		serviceID string
		want      string
	}{
		{
			name:      "nomad task service id",
			serviceID: "_nomad-task-" + allocID + "-web-http",
			want:      allocID,
		},
		{
			name:      "no uuid",
			serviceID: "postgres",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUUID(tt.serviceID); got != tt.want {
				t.Fatalf("extractUUID(%q) = %q, want %q", tt.serviceID, got, tt.want)
			}
		})
	}
}
