# CLAUDE.md — yeet-to-prod

Small, self-contained app. Read this before changing behavior — the domain rules and
architecture invariants below are the whole spec.

## What it is
Full-screen app answering "should I deploy right now?" — green YES / red NO with a
meme/movie-quote line, plus a live clock in the chosen country's timezone. A play on
"don't deploy on Friday." (API endpoint is `/should-i-deploy`.)

- `backend/` — Go API, one endpoint, **stdlib only** (no external deps — keeps `go test` /
  `go build` offline and the build reproducible).
- `frontend/` — Vite React, full-screen color + message + live country clock.

## Core rule (Safe to Deploy)
`safe = true` only when ALL hold, evaluated in the country's timezone:
**not night**, today **not weekend**, tomorrow **not weekend**, today **not holiday**,
tomorrow **not holiday**. Else `false`. (Rationale: never deploy when nobody's around to
fix a break.)

- **Night** = local time outside `work_hours` (before `start` or at/after `end`).
- Supported countries: `ID`, `IN`, `CN`, `US`, `AE` (Dubai/UAE), `JP` — all in `config.json`.

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
- `import _ "time/tzdata"` in `main.go` bundles tz data so `LoadLocation` works everywhere.

## Config (all via env / `.env`, see `.env.example`)
Backend: `PORT`, `ALLOWED_ORIGIN` (comma-separated list, or `*`), `CONFIG_PATH`,
`MESSAGES_PATH`, `DEFAULT_COUNTRY`, `RATE_LIMIT_PER_MINUTE`, `NOW_OVERRIDE` (fake clock,
RFC3339).
Frontend: `VITE_API_URL`, `VITE_DEFAULT_COUNTRY`.

## API
`GET /should-i-deploy?country=ID` → `{country, timezone, safe, reason, message}`; `400`
unknown country; `429` rate-limited. `GET /healthz` → `ok`. `timezone` drives the frontend
live clock (`formatClock` in `frontend/src/logic.js`).

## Security (honest)
CORS = browser hygiene only, not real auth (curl bypasses). `ALLOWED_ORIGIN` is a
comma-separated allowlist; the middleware echoes the request Origin when it matches (a
response can only name one origin), or sends `*` if the list contains `*`. In-memory rate
limit = single-instance only (resets on restart, not shared across replicas).

## Run
See `README.md`. Backend: `go test ./... && go run .`. Frontend: `npm install && npm run
dev`. Lint: `npm run lint`. Tests: `go test ./...`, `npm test`.

## Frontend deploy (GitHub Pages)
`.github/workflows/pages.yml`: push to any branch touching `frontend/**` → lint + test;
`main` → build + deploy to Pages (`deploy` needs `build` needs `test`). Build-time vars
`VITE_API_URL` and `VITE_DEFAULT_COUNTRY` come from repo **Variables** (Settings →
Variables), injected as env — never committed. Vite `base` = `/yeet-to-prod/` on `build`
only (dev stays `/`). One-time: set Pages Source = "GitHub Actions" + add the
`VITE_API_URL` and `VITE_DEFAULT_COUNTRY` Variables. Backend deploys separately (Render);
its `ALLOWED_ORIGIN` must include `https://qornanali.github.io`.

## Backend deploy (Render)
`.github/workflows/backend.yml`: push to any branch touching `backend/**` → `go vet` +
`go test`. Deploy is **manual only** (`workflow_dispatch` from `main`, `needs: test`) — a
human decision = Continuous Delivery (contrast: frontend Pages = auto = Continuous
Deployment). Deploy = `curl` of the Render **Deploy Hook** with `&ref=<github.sha>` so
Render builds the exact dispatched commit. Hook URL is a repo **Secret**
`RENDER_DEPLOY_HOOK`. Render Auto-Deploy must be OFF so this is the only deploy path.
