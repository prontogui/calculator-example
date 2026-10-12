package main

import (
	"testing"

	pg "github.com/prontogui/golib/v2"
)

// press finds the key whose Tag is code and delivers a Press to the
// calculator, the way Run does (UpdateState, then Refresh).
func press(t *testing.T, c *Calculator, keys map[string]*pg.Button, codes ...string) {
	t.Helper()
	for _, code := range codes {
		btn, ok := keys[code]
		if !ok {
			t.Fatalf("no key %q", code)
		}
		if !c.UpdateState(pg.Notification{Source: btn, Type: pg.NotifyPress}) {
			t.Fatalf("key %q not handled", code)
		}
		c.Refresh()
	}
}

// keysOf collects the calculator's buttons by Tag.
func keysOf(c *Calculator) map[string]*pg.Button {
	keys := map[string]*pg.Button{}
	var walk func(p pg.Primitive)
	walk = func(p pg.Primitive) {
		switch v := p.(type) {
		case *pg.Button:
			keys[v.Tag()] = v
		case *pg.Frame:
			for _, it := range v.Items() {
				walk(it)
			}
		case *pg.Group:
			for _, it := range v.Items() {
				walk(it)
			}
		}
	}
	walk(c.Content())
	return keys
}

// The GUI wiring around the engine (engine_test.go covers the math).
func Test_Calculator(t *testing.T) {
	c := NewCalculator()
	keys := keysOf(c)
	if len(keys) != 17 {
		t.Fatalf("found %d keys, want 17", len(keys))
	}

	// Every key's Tag must be a key the engine accepts.
	for code := range keys {
		if !NewEngine().Press(code) {
			t.Errorf("key Tag %q is not an engine key", code)
		}
	}

	// Presses reach the engine and Refresh shows its display.
	press(t, c, keys, "7", "*", "6", "=")
	if got := c.display.Content(); got != "42" {
		t.Errorf("display %q, want %q", got, "42")
	}
	if got := c.display.Content(); got != c.engine.Display() {
		t.Errorf("display %q differs from the engine's %q", got, c.engine.Display())
	}

	// A non-press notification is left for someone else.
	if c.UpdateState(pg.Notification{Source: keys["1"], Type: pg.NotifyRelease}) {
		t.Error("a Release must not be handled")
	}
}
