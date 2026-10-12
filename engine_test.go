package main

import "testing"

// The engine on its own - no GUI involved.
func Test_Engine(t *testing.T) {
	cases := []struct {
		keys string // one key per character
		want string
	}{
		{"", "0"},
		{"12+3=", "15"},
		{"7*6=", "42"},
		{"9/0=", "0"},    // division by zero gives 0
		{"1.5+.5=", "2"}, // a leading "." starts "0."
		{"2+3*4=", "20"}, // left to right, like a pocket calculator
		{"5+", "5"},      // an operator keeps showing the operand
		{"5+3", "3"},     // the next digit starts a new number
		{"5+3+", "8"},    // a second operator completes the first
		{"1..2", "1.2"},  // a second "." is ignored
		{"007", "7"},     // leading zeros are dropped
		{"8-10=", "-2"},
		{"12+3=C", "0"}, // clear
		{"2+3==", "5"},  // "=" with nothing pending does nothing
		{"2+3=4", "4"},  // a digit after "=" starts over
	}
	for _, tc := range cases {
		e := NewEngine()
		for _, k := range tc.keys {
			if !e.Press(string(k)) {
				t.Fatalf("%q: key %q rejected", tc.keys, k)
			}
		}
		if got := e.Display(); got != tc.want {
			t.Errorf("%q: display %q, want %q", tc.keys, got, tc.want)
		}
	}

	if NewEngine().Press("x") {
		t.Error("an unknown key must be rejected")
	}
}
