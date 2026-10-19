package main

import "testing"

func TestBuildDiagnostic(t *testing.T) {
	tests := []struct{ name, status, code string }{{"passed", "passed", "BUILD_OK"}, {"failed", "failed", "BUILD_FAILED"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDiagnostic(BuildEvent{Commit: "abc123", Status: tt.status})
			if got.Code != tt.code {
				t.Fatalf("code=%q want %q", got.Code, tt.code)
			}
		})
	}
}
