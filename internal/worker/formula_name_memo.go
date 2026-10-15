package worker

import "sync"

// formulaNameMemoCap bounds the memo. Past the cap the memo is cleared rather
// than evicted piecemeal, which costs one re-read per live acting bead and
// keeps the structure trivial.
const formulaNameMemoCap = 1024

// FormulaNameMemo remembers the formula an acting bead's run was compiled
// from, keyed by the acting bead id.
type FormulaNameMemo struct {
	mu      sync.Mutex
	entries map[string]string
}

// NewFormulaNameMemo returns an empty memo.
func NewFormulaNameMemo() *FormulaNameMemo {
	return &FormulaNameMemo{entries: make(map[string]string)}
}

// resolve runs derive; the skeleton does not memoize yet.
func (m *FormulaNameMemo) resolve(_ string, derive func() (string, error)) (string, error) {
	return derive()
}
