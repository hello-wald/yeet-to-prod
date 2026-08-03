# yeet-to-prod

**Should you deploy right now?** This little app gives you a brutally honest answer. It fills
your whole screen with **green (YES)** or **red (NO)**, adds a cheeky one-liner (borrowed from
movies, games, and memes), and shows a live clock for the country you pick.

It's a friendly joke about a very real engineering superstition: *don't deploy on a Friday
(or a holiday, or the middle of the night) — because if something breaks, nobody's around to
fix it.*

The project has two halves:

- **`backend/`** — a small API written in Go. It has a single endpoint and a full set of
  tests, and it uses only Go's standard library (no third-party packages to install).
- **`frontend/`** — a React app built with Vite. It shows the full-screen verdict and the
  live country clock.

## How it decides

The app says **YES** only when *all* of these are true, checked in the selected country's
timezone:

- it's **not night** (you're within normal working hours),
- **today** is not a weekend,
- **tomorrow** is not a weekend,
- **today** is not a holiday, and
- **tomorrow** is not a holiday.

If any one of those fails, the answer is **NO**. (The idea: only ship when someone will be
around tomorrow to handle a problem.)

Everything the decision depends on lives in plain JSON, so you can tweak it without touching
code:

- **`backend/config.json`** — the countries (`ID`, `IN`, `CN`, `US`, `AE`, `JP`) and their
  timezones, weekends, working hours, and holidays.
- **`backend/reasons.json`** — the pool of funny messages. The app picks one at random. Add
  or remove lines here anytime.

## Running it locally

### Backend (requires Go 1.22 or newer)

```bash
cd backend
cp .env.example .env          # optional — the defaults work fine
go test ./...                 # run the full test suite
go run .                      # starts the API on http://localhost:8080
```

Once it's running, try a couple of requests:

```bash
curl "http://localhost:8080/should-i-deploy?country=ID"
curl "http://localhost:8080/should-i-deploy?country=AE"
```

### Frontend (requires Node 20 or newer)

```bash
cd frontend
cp .env.example .env          # set VITE_API_URL to point at your backend
npm install                   # installs dependencies (creates package-lock.json)
npm run lint
npm test
npm run dev                   # opens the app on http://localhost:5173
```

### Testing any day or time

You don't have to wait until Friday night to see a "NO". Set a fake clock in the backend
`.env` file:

```
NOW_OVERRIDE=2026-08-15T10:00:00+07:00   # a Saturday in Jakarta → NO
```

Leave it empty to use the real current time.

## The API

There's just one main endpoint:

```
GET /should-i-deploy?country=ID
→ { "country": "ID", "timezone": "Asia/Jakarta", "safe": false,
    "reason": "it's the weekend", "message": "..." }
```

An unknown country returns `400`, too many requests return `429`, and `GET /healthz`
returns `ok` (handy for uptime checks).

**A note on access control.** The backend uses CORS to restrict which *websites* (origins)
may call it, configured via `ALLOWED_ORIGIN` (a comma-separated list, or `*` for any). Keep
in mind CORS is enforced by the browser only — it won't stop tools like `curl` or Postman,
so treat it as good hygiene rather than real security. Likewise, the built-in rate limit
lives in memory, so it resets when the server restarts and isn't shared across multiple
instances. Both are perfectly fine for a small public app like this one.

## Deploying (all on free tiers)

The repo ships with two GitHub Actions workflows. Together they show the two flavours of CD:
the frontend deploys itself automatically (**Continuous Deployment**), while the backend
waits for you to press a button (**Continuous Delivery**).

### Frontend → GitHub Pages (automatic)

`.github/workflows/pages.yml` runs lint and tests on every push that touches `frontend/`. On
the `main` branch it also builds the site and deploys it to GitHub Pages. The jobs are
chained so that `deploy` only runs if `build` succeeds, which only runs if `test` passes.

To set it up once:

1. Go to **Settings → Pages** and set **Source** to **"GitHub Actions"**.
2. Go to **Settings → Secrets and variables → Actions → Variables** and add
   `VITE_API_URL` with your backend's URL. Use a *Variable* (not a Secret) — it's a public
   URL, and the workflow bakes it into the build, so nothing sensitive is committed.
3. Your site will be live at `https://qornanali.github.io/yeet-to-prod/`. (Vite's `base` is
   already set to match this path.)

Two things to remember: GitHub Pages is HTTPS-only, so your backend must also be HTTPS, and
you'll need to add the Pages address (`https://qornanali.github.io`) to the backend's
`ALLOWED_ORIGIN` so the browser is allowed to call it.

### Backend → Render (manual, gated by tests)

`.github/workflows/backend.yml` runs `go vet` and the tests on every push that touches
`backend/`. Deploying, however, is a deliberate manual step: go to the **Actions** tab,
choose **Run workflow** on `main`, and once the tests pass it tells Render to deploy the
exact commit you ran it on. That human "go" is what makes it Continuous Delivery.

To set it up once:

1. Create a Render web service with **Root Directory** = `backend/`, build command
   `go build -o app .`, and start command `./app`.
2. On Render, set `ALLOWED_ORIGIN` to include `https://qornanali.github.io`.
3. On Render, **turn off Auto-Deploy** — that way GitHub Actions is the only thing that
   deploys, so your tests always act as the gate.
4. On Render, copy the **Deploy Hook** URL and add it to GitHub as a **Secret** named
   `RENDER_DEPLOY_HOOK`. (A Secret, not a Variable, because anyone who has this URL can
   trigger a production deploy.)

One quirk of the free tier: the backend goes to sleep after about 15 minutes of inactivity,
so the first request after a nap takes an extra ~30 seconds to wake it up.
