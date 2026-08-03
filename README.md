# yeet-to-prod

A full-screen app that gives a brutally honest answer to: **should you deploy right now?**
Green YES / red NO with a meme/movie-quote line, plus a live clock in the chosen country's
timezone. A play on the "don't deploy on Friday" rule. (API endpoint is `/should-i-deploy`.)

- **backend/** — Go API, one endpoint, pure `decide()` + full test suite, stdlib only.
- **frontend/** — Vite React, full-screen green/red verdict + live country clock.

Rule (**Safe to Deploy**): YES only when it's **not night**, and **today + tomorrow** are
neither **weekend** nor **holiday** (per country). Anything else → NO.

Countries, holidays, work hours, weekend days → `backend/config.json` (`ID/IN/CN/US/AE/JP`).
Messages (randomized) → `backend/reasons.json`. Add/remove = edit JSON.

---

## Run locally

### Backend (needs Go 1.22+)
```bash
cd backend
cp .env.example .env          # optional; defaults work
go test ./...                 # run the full test suite
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

## Fake the clock (test any day/time)

Backend `.env`:
```
NOW_OVERRIDE=2026-08-15T10:00:00+07:00   # a Saturday in Jakarta → NO
```
Empty = real time.

---

## API

`GET /should-i-deploy?country=ID` → `{ country, timezone, safe, reason, message }`
- `400` unknown country · `429` rate-limited · `GET /healthz` → `ok`

**Note on access control:** CORS restricts *browsers* to the origins in `ALLOWED_ORIGIN`
(comma-separated list, or `*`) but does not stop `curl`/Postman — it's hygiene, not auth.
The in-memory rate limit is per-instance (resets on restart, not shared across replicas).
Both are fine for a small public app.

---

## Deploy (free tier)

- **frontend/** → Vercel · Root Directory = `frontend/` · set `VITE_API_URL` to the backend URL.
- **backend/** → Render (web service) · Root Directory = `backend/` · build `go build -o app .`,
  start `./app`. Set `ALLOWED_ORIGIN` to the frontend URL. Free instances cold-start after
  ~15 min idle.
