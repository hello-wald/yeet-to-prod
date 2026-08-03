package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testMessages() Messages {
	return Messages{
		CodeOK:           {Reason: "all clear", Messages: []string{"YES go"}},
		CodeTodayWeekend: {Reason: "it's the weekend", Messages: []string{"NO weekend"}},
		CodeNight:        {Reason: "it's night", Messages: []string{"NO night"}},
	}
}

// newTestServer builds a server with in-code config/messages and a generous
// rate limit. nowOverride pins the clock so decide() is deterministic.
func newTestServer(nowOverride string, maxPerMin int) *server {
	return &server{
		cfg:            Config{"ID": idCfg, "SA": saCfg},
		msgs:           testMessages(),
		lim:            newLimiter(maxPerMin),
		rng:            rand.New(rand.NewSource(1)),
		defaultCountry: "ID",
		nowOverride:    nowOverride,
	}
}

func do(s *server, method, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	s.router([]string{"http://localhost:5173"}).ServeHTTP(rec, req)
	return rec
}

// doOrigin sends a request with an Origin header against a given allowlist.
func doOrigin(s *server, target, origin string, allowed []string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", target, nil)
	req.Header.Set("Origin", origin)
	rec := httptest.NewRecorder()
	s.router(allowed).ServeHTTP(rec, req)
	return rec
}

func TestHandleDeploy_Safe(t *testing.T) {
	// 2026-08-12 Wed 10:00 Jakarta → safe
	rec := do(newTestServer("2026-08-12T10:00:00+07:00", 100), "GET", "/should-i-deploy?country=ID")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var r response
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if !r.Safe || r.Reason != "all clear" || r.Country != "ID" || r.Message == "" {
		t.Errorf("got %+v; want safe=true reason='all clear' country=ID non-empty message", r)
	}
}

func TestHandleDeploy_Unsafe(t *testing.T) {
	// 2026-08-15 Sat → weekend
	rec := do(newTestServer("2026-08-15T10:00:00+07:00", 100), "GET", "/should-i-deploy?country=ID")
	var r response
	_ = json.Unmarshal(rec.Body.Bytes(), &r)
	if r.Safe || r.Reason != "it's the weekend" {
		t.Errorf("got %+v; want safe=false reason=\"it's the weekend\"", r)
	}
}

func TestHandleDeploy_DefaultCountry(t *testing.T) {
	// no country param → defaults to ID
	rec := do(newTestServer("2026-08-12T10:00:00+07:00", 100), "GET", "/should-i-deploy")
	var r response
	_ = json.Unmarshal(rec.Body.Bytes(), &r)
	if r.Country != "ID" {
		t.Errorf("country = %q, want ID (default)", r.Country)
	}
}

func TestHandleDeploy_UnknownCountry(t *testing.T) {
	rec := do(newTestServer("2026-08-12T10:00:00+07:00", 100), "GET", "/should-i-deploy?country=ZZ")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestHandleDeploy_RateLimited(t *testing.T) {
	s := newTestServer("2026-08-12T10:00:00+07:00", 1)
	if r := do(s, "GET", "/should-i-deploy?country=ID"); r.Code != http.StatusOK {
		t.Fatalf("first call status = %d, want 200", r.Code)
	}
	if r := do(s, "GET", "/should-i-deploy?country=ID"); r.Code != http.StatusTooManyRequests {
		t.Errorf("second call status = %d, want 429", r.Code)
	}
}

func TestHealthz(t *testing.T) {
	rec := do(newTestServer("", 100), "GET", "/healthz")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Errorf("healthz = %d %q; want 200 ok", rec.Code, rec.Body.String())
	}
}

func TestCORS(t *testing.T) {
	s := newTestServer("2026-08-12T10:00:00+07:00", 100)
	allow := []string{"http://localhost:5173", "https://yeet.example.com"}

	// each allowed origin is echoed back (that's how a multi-origin allowlist works)
	for _, o := range allow {
		rec := doOrigin(s, "/should-i-deploy?country=ID", o, allow)
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != o {
			t.Errorf("origin %q: ACAO = %q, want it echoed", o, got)
		}
	}

	// disallowed origin → no ACAO header at all
	rec := doOrigin(s, "/should-i-deploy?country=ID", "https://evil.example.com", allow)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("disallowed origin got ACAO %q, want empty", got)
	}

	// "*" allows any origin
	star := doOrigin(s, "/should-i-deploy?country=ID", "https://anything.example.com", []string{"*"})
	if got := star.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("wildcard ACAO = %q, want *", got)
	}

	// preflight OPTIONS → 204
	pre := do(newTestServer("", 100), "OPTIONS", "/should-i-deploy")
	if pre.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, want 204", pre.Code)
	}
}
