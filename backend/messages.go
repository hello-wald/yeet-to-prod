package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// Verdict is the human-facing text for one reason code: a factual reason plus a
// pool of funny lines to pick from at random.
type Verdict struct {
	Reason   string   `json:"reason"`
	Messages []string `json:"messages"`
}

// Messages maps reason code → Verdict. Loaded from reasons.json so lines can be
// added/removed/edited without touching Go code.
type Messages map[string]Verdict

func loadMessages(path string) (Messages, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read messages %s: %w", path, err)
	}
	var m Messages
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parse messages %s: %w", path, err)
	}
	if _, ok := m[CodeOK]; !ok {
		return nil, fmt.Errorf("messages %s missing required %q entry", path, CodeOK)
	}
	return m, nil
}

// lookup returns the Verdict for a code, falling back safely if the file lacks it.
func (m Messages) lookup(code string) Verdict {
	if v, ok := m[code]; ok {
		return v
	}
	return Verdict{Reason: code, Messages: []string{"NO 🙅 not today."}}
}

// pick returns a random line from the Verdict's pool. rng is injected so the
// random source stays swappable/testable.
func (v Verdict) pick(rng *rand.Rand) string {
	if len(v.Messages) == 0 {
		return "NO 🙅 not today."
	}
	return v.Messages[rng.Intn(len(v.Messages))]
}
