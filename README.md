# yeet-to-prod

Fun demo for the COMPFEST 18 CI/CD session. A full-screen app that tells you whether it's
safe to deploy right now — should you yeet to prod, or not? Meta-joke on the "don't deploy
on Friday" rule. (API endpoint stays `/should-i-deploy`.)

- **backend/** — Go API, one endpoint, pure `decide()` + table tests, stdlib only.
- **frontend/** — Vite React, full-screen green/red + funny line.

Rule (**Safe to Deploy**): YES only when it's **not night**, and **today + tomorrow** are
neither **weekend** nor **holiday** (per country). Anything else → NO.

Countries, holidays, work hours, weekend days → `backend/config.json`.
Funny messages (randomized) → `backend/reasons.json`. Add/remove = edit JSON.

---

## Run locally

### Backend (needs Go 1.22+)
```bash
cd backend
cp .env.example .env          # optional; defaults work
go test ./...                 # run the tests CI will run
go run .                      # serves http://localhost:8080
```
Try it:
```bash
curl "http://localhost:8080/should-i-deploy?country=ID"
curl "http://localhost:8080/should-i-deploy?country=AE"
```

### Frontend (needs Node 20+)
```bash
cd frontend
cp .env.example .env          # point VITE_API_URL at the backend
npm install                   # creates package-lock.json
npm run lint
npm test
npm run dev                   # http://localhost:5173
```

---

## Fake the clock (demo without waiting)

Backend `.env`:
```
NOW_OVERRIDE=2026-08-15T10:00:00+07:00   # a Saturday in Jakarta → NO
```
Empty = real time.

---

## The two on-stage breaks (for the CI/CD lesson)

**1. Logic break (backend test goes red).** In `backend/decide.go`, invert a weekend check:
```go
// green:
if isWeekend(tomorrow.Weekday(), cfg.WeekendDays) {
// broken — flip the condition:
if !isWeekend(tomorrow.Weekday(), cfg.WeekendDays) {
```
Now the tool says "deploy on Saturday." `go test ./...` fails on
`tomorrow weekend blocks`. Revert → green. Teaches: `deploy needs: test`.

**2. Dependency break (frontend `npm ci` goes red).** Delete `frontend/package-lock.json`
and commit. On CI, `npm ci` fails ("works on my machine"). Commit the lockfile → green.
Teaches: dependency management, reproducible builds.

---

## Deploy (free tier)

- **frontend/** → Vercel · Root Directory = `frontend/` · set `VITE_API_URL` to the backend URL.
- **backend/** → Render (web service) · Root Directory = `backend/` · build `go build -o app .`,
  start `./app`. Set `ALLOWED_ORIGIN` to the Vercel URL. Pre-warm before demo (free instance
  cold-starts after 15 min idle).
