// A four-function calculator built with ProntoGUI 2.0.
//
// The math lives in Engine (engine.go), which knows nothing about the GUI.
// This file is the GUI around it, and shows:
//   - a Component (Calculator) that handles its own keys: UpdateState passes
//     a key press to the engine, Refresh shows the result;
//   - pgui.Run, which replaces a hand-written SetGUI + Wait loop;
//   - notifications via pg.EnableNotify plus a Tag per key, instead of a
//     closure per key;
//   - optional layout fields, set with Go 1.26 new(...), e.g. new(16.0).
package main

import (
	"fmt"

	pg "github.com/prontogui/golib/v2"
)

// Calculator is a reusable component: a display and a keypad. Embedding
// pg.ComponentBase makes it a primitive you can place anywhere in a GUI;
// UpdateState and Refresh make it handle its own keys.
type Calculator struct {
	pg.ComponentBase

	engine  *Engine
	display *pg.Text
}

// NewCalculator builds the calculator's primitives.
func NewCalculator() *Calculator {
	c := &Calculator{engine: NewEngine()}

	c.display = pg.TextWith{
		FontFamily:              "Roboto",
		FontSize:                36.0,
		HorizontalAlignment:     pg.HorizontalAlignmentExpand,
		HorizontalTextAlignment: pg.HorizontalTextAlignmentRight,
	}.Make()

	// Optional fields: pointers written with new(...). Only what's set is
	// sent; anything left out uses the App's default.
	displayFrame := pg.FrameWith{
		Items:               pg.VA(c.display),
		PaddingAll:          new(16.0),
		MarginBottom:        new(8.0),
		BorderAll:           new(1.0),
		HorizontalAlignment: pg.HorizontalAlignmentExpand,
	}.Make()

	// Every key reports its press (EnableNotify) and carries its engine key
	// code in Tag; UpdateState below passes them all to the engine.
	key := func(label, code string, style pg.ButtonStyle, width float64) *pg.Button {
		return pg.ButtonWith{
			Label:   label,
			Tag:     code,
			Style:   style,
			Width:   width,
			Height:  56.0,
			OnPress: pg.EnableNotify,
		}.Make()
	}
	k := func(label, code string) *pg.Button { return key(label, code, pg.ButtonStyleElevated, 64.0) }
	row := func(items ...pg.Primitive) *pg.Group { return pg.GroupWith{Items: items}.Make() }

	keypad := pg.FrameWith{
		Items: pg.VA(
			row(k("7", "7"), k("8", "8"), k("9", "9"), k("÷", "/")),
			row(k("4", "4"), k("5", "5"), k("6", "6"), k("×", "*")),
			row(k("1", "1"), k("2", "2"), k("3", "3"), k("−", "-")),
			row(k("C", "C"), k("0", "0"), k(".", "."), k("+", "+")),
			row(key("=", "=", pg.ButtonStyleOutlined, 280.0)),
		),
		HorizontalAlignment: pg.HorizontalAlignmentCenter,
	}.Make()

	c.SetContent(pg.FrameWith{Items: pg.VA(displayFrame, keypad)}.Make())
	c.Refresh()
	return c
}

// UpdateState passes a key press inside the calculator to the engine.
// Returning true marks the notification handled, so it doesn't bubble any
// further.
func (c *Calculator) UpdateState(n pg.Notification) bool {
	btn, ok := n.Source.(*pg.Button)
	if !ok || n.Type != pg.NotifyPress {
		return false
	}
	return c.engine.Press(btn.Tag())
}

// Refresh shows the engine's display on screen.
func (c *Calculator) Refresh() {
	c.display.SetContent(c.engine.Display())
}

func main() {
	pgui := pg.NewProntoGUI()
	if err := pgui.StartServingSingle("127.0.0.1", 50053); err != nil {
		fmt.Printf("Error starting server: %s\n", err.Error())
		return
	}
	defer pgui.StopServing()

	// FlowDirection is optional and left unset, so the View uses its default
	// (top to bottom).
	view := pg.ViewWith{
		Items:   pg.VA(NewCalculator()),
		Title:   "Calculator",
		Showing: true,
	}.Make()

	// Run sends the GUI and delivers each key press to the Calculator. It
	// handles the App disconnecting and reconnecting, and returns only on an
	// error.
	if err := pgui.Run(pg.RunWithArgs{Primitives: pg.VA(view)}); err != nil {
		fmt.Printf("Run error: %s\n", err.Error())
	}
}
