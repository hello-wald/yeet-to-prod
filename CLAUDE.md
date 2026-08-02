# CLAUDE.md — yeet-to-prod

Teaching demo for the COMPFEST 18 CI/CD session. Small on purpose. Full domain design +
glossary lives in `../CONTEXT.md` — read it before changing behavior.

## What it is
Full-screen app answering "should I yeet to prod right now?" — green YES / red NO with a
meme/movie-quote line, plus a live clock in the chosen country's timezone. Meta-joke on
"don't deploy on Friday." (Brand = yeet-to-prod; API endpoint stays `/should-i-deploy`.)

- `backend/` — Go API, one endpoint, **stdlib only** (no external deps — keeps `go test` /
  `go build` offline + the dependency story clean).
- `frontend/` — Vite React, full-screen color + message + live country clock.

## Core rule (Safe to Deploy)
`true` only when ALL hold: **not night**, today **not weekend**, tomorrow **not weekend**,
today **not holiday**, tomorrow **not holiday**. Else `false`.

Supported countries: `ID`, `IN`, `CN`, `US`, `AE` (Dubai/UAE), `JP`. All in `config.json`.

## Architecture invariants — keep these
- **`decide(now, cfg)` is PURE.** Time is **injected** (never `time.Now()` inside) so tests
  pin any clock. Don't add I/O, globals, or `rand` here.
- **`decide` returns a stable reason CODE**, not prose (`ok`, `night`, `today_weekend`,
  `tomorrow_weekend`, `today_holiday`, `tomorrow_holiday`). Check order = reason priority.
- **Prose is data, not code.** `reasons.json` maps code → `{reason, messages[]}`. `message`
  = random pick from the array. Add/remove lines = edit JSON only.
- **Country data is data.** `config.json` maps country ID → `{timezone, weekend_days,
  work_hours, holidays[]}`. Schema supports per-country weekends (`[6,0]` Sat/Sun,
  `[5,6]` Fri/Sat) — all current countries are `[6,0]`.
- **Randomness isolated** in the message picker (injected `*rand.Rand`), not in `decide`.
- **HTTP handler lives in `server.go`** (`server` type), split from `main.go` wiring so it's
  testable with `httptest`. `main.go` only loads config + env and starts the listener.
- `import _ "time/tzdata"` in `main.go` bundles tz data so `LoadLocation` works on CI/Render.

## Config (all via env / `.env`, see `.env.example`)
Backend: `PORT`, `ALLOWED_ORIGIN`, `CONFIG_PATH`, `MESSAGES_PATH`, `DEFAULT_COUNTRY`,
`RATE_LIMIT_PER_MINUTE`, `NOW_OVERRIDE` (fake clock, RFC3339, demo only).
Frontend: `VITE_API_URL`, `VITE_DEFAULT_COUNTRY`.

## API
`GET /should-i-deploy?country=ID` → `{country, timezone, safe, reason, message}`; `400`
unknown country; `429` rate-limited. `GET /healthz` → `ok`. `timezone` drives the frontend
live clock (`formatClock` in `frontend/src/logic.js`).

## Security (honest)
CORS = browser hygiene only, not real auth (curl bypasses). In-memory rate limit =
single-instance only (resets on restart, not shared across replicas). Fine for a demo;
say so out loud — it's a teaching point.

## Run
See `README.md`. Backend: `go test ./... && go run .`. Frontend: `npm install && npm run
dev`. Lint: `npm run lint`. Tests: `go test ./...`, `npm test`.

## The two on-stage breaks (don't "fix" these — they're the lesson)
1. Invert a weekend check in `decide.go` → backend test red → teaches `deploy needs: test`.
2. Delete `frontend/package-lock.json` → `npm ci` red → teaches dependency management.
See `README.md` for exact steps.
