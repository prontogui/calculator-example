package main

import (
	"strconv"
	"strings"
)

// Engine is the calculator's model: its state and arithmetic. It knows
// nothing about the GUI - the Calculator component feeds it key presses and
// shows Display().
//
// Operations are evaluated left to right as they are entered (2 + 3 × 4 is
// 20), like a basic pocket calculator.
type Engine struct {
	display    string  // what the calculator shows
	stored     float64 // left operand of the pending operation
	pendingOp  string  // "+", "-", "*", "/", or "" for none
	resetOnKey bool    // the next digit starts a new number
}

// NewEngine returns an engine showing 0.
func NewEngine() *Engine {
	return &Engine{display: "0"}
}

// Display returns what the calculator should show.
func (e *Engine) Display() string {
	return e.display
}

// Press applies one key: "0"-"9", ".", "+", "-", "*", "/", "=", or "C"
// (clear). It reports false for any other key.
func (e *Engine) Press(key string) bool {
	switch key {
	case "C":
		e.Clear()
	case "=":
		e.Equals()
	case ".":
		e.Decimal()
	case "+", "-", "*", "/":
		e.Operator(key)
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		e.Digit(key)
	default:
		return false
	}
	return true
}

// Digit appends a digit, or starts a new number after an operator or "=".
func (e *Engine) Digit(d string) {
	if e.resetOnKey || e.display == "0" {
		e.display = d
		e.resetOnKey = false
	} else {
		e.display += d
	}
}

// Decimal adds a decimal point, if the number doesn't have one.
func (e *Engine) Decimal() {
	if e.resetOnKey {
		e.display = "0."
		e.resetOnKey = false
		return
	}
	if !strings.Contains(e.display, ".") {
		e.display += "."
	}
}

// Operator starts op ("+", "-", "*", "/"), first completing any pending
// operation.
func (e *Engine) Operator(op string) {
	val, _ := strconv.ParseFloat(e.display, 64)
	if e.pendingOp != "" && !e.resetOnKey {
		e.stored = compute(e.stored, val, e.pendingOp)
		e.display = formatResult(e.stored)
	} else {
		e.stored = val
	}
	e.pendingOp = op
	e.resetOnKey = true
}

// Equals completes the pending operation, if any.
func (e *Engine) Equals() {
	if e.pendingOp == "" {
		return
	}
	val, _ := strconv.ParseFloat(e.display, 64)
	e.stored = compute(e.stored, val, e.pendingOp)
	e.display = formatResult(e.stored)
	e.pendingOp = ""
	e.resetOnKey = true
}

// Clear resets the engine to 0.
func (e *Engine) Clear() {
	*e = Engine{display: "0"}
}

// compute applies op to a and b. Division by zero gives 0.
func compute(a, b float64, op string) float64 {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			return 0
		}
		return a / b
	}
	return b
}

func formatResult(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}
