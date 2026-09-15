package repository

import "testing"

func TestNormalizePointType(t *testing.T) {
	cases := map[string]string{
		"credit":  "credit",
		"earned":  "credit",
		"debit":   "debit",
		"spent":   "debit",
		"unknown": "unknown",
	}

	for input, expected := range cases {
		if got := NormalizePointType(input); got != expected {
			t.Fatalf("NormalizePointType(%q) = %q, want %q", input, got, expected)
		}
	}
}
