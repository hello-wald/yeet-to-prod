package main

import (
	"math/rand"
	"testing"
)

func TestLoadMessages_OK(t *testing.T) {
	p := writeTemp(t, "reasons.json", `{
		"ok": {"reason":"all clear","messages":["a","b"]},
		"night": {"reason":"it's night","messages":["z"]}
	}`)
	m, err := loadMessages(p)
	if err != nil {
		t.Fatal(err)
	}
	if m[CodeOK].Reason != "all clear" || len(m[CodeOK].Messages) != 2 {
		t.Errorf("parsed messages wrong: %+v", m)
	}
}

func TestLoadMessages_MissingOK(t *testing.T) {
	p := writeTemp(t, "reasons.json", `{"night":{"reason":"x","messages":["y"]}}`)
	if _, err := loadMessages(p); err == nil {
		t.Error("want error when 'ok' entry missing, got nil")
	}
}

func TestLoadMessages_MissingFile(t *testing.T) {
	if _, err := loadMessages("nope.json"); err == nil {
		t.Error("want error for missing file, got nil")
	}
}

func TestLookup_Fallback(t *testing.T) {
	m := Messages{CodeOK: {Reason: "all clear", Messages: []string{"go"}}}
	v := m.lookup("unknown_code")
	if v.Reason != "unknown_code" || len(v.Messages) == 0 {
		t.Errorf("fallback wrong: %+v", v)
	}
}

func TestPick_InPool(t *testing.T) {
	v := Verdict{Messages: []string{"a", "b", "c"}}
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 20; i++ {
		got := v.pick(rng)
		if got != "a" && got != "b" && got != "c" {
			t.Fatalf("pick returned %q, not in pool", got)
		}
	}
}

func TestPick_EmptyPool(t *testing.T) {
	v := Verdict{Messages: nil}
	if got := v.pick(rand.New(rand.NewSource(1))); got == "" {
		t.Error("pick on empty pool returned empty string, want fallback")
	}
}
