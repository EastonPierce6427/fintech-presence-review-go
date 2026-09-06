package main

import "testing"

func TestDecideRisk(t *testing.T) {
	tests := []struct {
		name    string
		event   PaymentEvent
		online  map[string]bool
		allowed bool
	}{
		{"online reviewer", PaymentEvent{Actor: "a", AmountCents: 5000}, map[string]bool{"a": true}, true},
		{"offline actor", PaymentEvent{Actor: "a", AmountCents: 5000}, map[string]bool{}, false},
		{"high value", PaymentEvent{Actor: "a", AmountCents: 100001}, map[string]bool{"a": true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecideRisk(tt.event, tt.online).Allowed; got != tt.allowed {
				t.Fatalf("allowed=%v, want %v", got, tt.allowed)
			}
		})
	}
}
