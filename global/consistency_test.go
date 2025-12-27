package global

import "testing"

func TestStrictModeEpochMismatch(t *testing.T) {
	err := CheckEpochAlignment([]uint64{1, 2}, Strict)
	if err == nil {
		t.Fatalf("expected epoch mismatch error")
	}
}

func TestFastModeAllowsMismatch(t *testing.T) {
	err := CheckEpochAlignment([]uint64{1, 2}, Fast)
	if err != nil {
		t.Fatalf("fast mode should allow mismatch")
	}
}
