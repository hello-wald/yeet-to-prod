package main

import (
	"testing"
	"time"
)

func TestLimiter_AllowsUnderMax(t *testing.T) {
	l := newLimiter(3)
	now := time.Now()
	for i := 0; i < 3; i++ {
		if !l.allow("1.1.1.1", now) {
			t.Fatalf("request %d denied, want allowed", i+1)
		}
	}
	if l.allow("1.1.1.1", now) {
		t.Error("4th request allowed, want denied (over max)")
	}
}

func TestLimiter_PerIP(t *testing.T) {
	l := newLimiter(1)
	now := time.Now()
	if !l.allow("1.1.1.1", now) || !l.allow("2.2.2.2", now) {
		t.Error("different IPs should each get their own budget")
	}
	if l.allow("1.1.1.1", now) {
		t.Error("same IP over max should be denied")
	}
}

func TestLimiter_WindowExpiry(t *testing.T) {
	l := newLimiter(1)
	base := time.Now()
	if !l.allow("1.1.1.1", base) {
		t.Fatal("first request denied")
	}
	if l.allow("1.1.1.1", base.Add(30*time.Second)) {
		t.Error("within window should still be denied")
	}
	if !l.allow("1.1.1.1", base.Add(61*time.Second)) {
		t.Error("after window should be allowed again")
	}
}
