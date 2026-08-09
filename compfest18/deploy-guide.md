# Deploy Guide — Frontend to GitHub Pages, Backend to Render

A follow-along guide to ship a small full-stack app for **free**, the CI/CD way:
you `git push`, and robots test + deploy your code automatically.

We use the demo app **`yeet-to-prod`** as the concrete example. Every value in a
`code block` that starts with `yeet-to-prod` or a URL is **specific to the demo** —
wherever you see **`← replace with yours`**, swap in your own project's value.

> You do **not** need Node or Go installed on your laptop. The robots (GitHub Actions)
> build everything in the cloud. You only need a browser and `git`.

### 📝 Your values — fill these in first

Every concrete value below (like `qornanali.github.io` or `yeet-to-prod.onrender.com`) is
**from the demo**. Wherever you see one, use your own from this table. Fill it as you go —
you'll know each value once you reach its step.

| What | Demo value (example) | **Yours** |
|---|---|---|
| GitHub username | `qornanali` | `__________` |
| Repo name | `yeet-to-prod` | `__________` |
| Frontend URL (Pages) | `https://qornanali.github.io/yeet-to-prod/` | `__________` |
| Backend URL (Render) | `https://yeet-to-prod.onrender.com` | `__________` |
| CORS origin (host only) | `https://qornanali.github.io` | `__________` |

> Tip: if you forked and kept the name `yeet-to-prod`, only the **username** part changes.

---

## 1. What you'll build

Two pieces, deployed to two free hosts, wired together:

- **Frontend** (static React site) → **GitHub Pages** — deploys **automatically** on every
  push to `main` (Continuous Deployment).
- **Backend** (Go API) → **Render** — deploys **when you click a button** (Continuous
  Delivery).

```
                    ┌─────────────┐
    you  ──push──▶  │ GitHub repo │
                    │ frontend/   │
                    │ backend/    │
                    └──────┬──────┘
                           │ triggers GitHub Actions (the robots)
              ┌────────────┴────────────┐
              ▼                          ▼
     pages.yml (auto)          backend.yml (manual)
     test → build → deploy      test → deploy (you click)
              │                          │
              ▼                          ▼
    ┌───────────────────┐      ┌───────────────────┐
    │   GitHub Pages     │      │      Render        │
    │ (static frontend)  │      │  (Go backend/API)  │
    └─────────┬──────────┘      └─────────▲──────────┘
              │  browser loads the page    │
              ▼                            │
        ┌───────────┐   fetch /should-i-deploy  (CORS)
        │  Browser  │ ───────────────────────────┘
        └───────────┘
```

**The flow:** you push code → GitHub Actions runs your tests → if green, it deploys. The
browser loads the frontend from Pages, and the frontend calls the backend on Render.

---

## 2. Prerequisites

1. **GitHub account** + your project pushed to a repo. Monorepo layout:
   ```
   your-repo/
     frontend/     # the React app
     backend/      # the Go API
   ```
2. **Render account** — free. Sign up at <https://render.com> using "Continue with GitHub".
3. **Git basics** — `git add`, `git commit`, `git push` (from the Collaborative Programming
   module).
4. **A web browser** — most steps here are clicks on GitHub and Render settings pages.

> Not required: Node.js or Go installed locally. CI builds in the cloud.

---

## 2.5. Step 0 — get the code

Pick one:

**Option A — Fork + start from the `starter` branch (recommended for practice).** Browser →
the demo repo → **Fork** (top-right) → creates `https://github.com/<your-user>/yeet-to-prod`.

The repo has two relevant branches:
- **`main`** — the **finished** result (workflows already added). Read it to compare answers.
- **`starter`** — the app **without** any CI/CD. This is where you follow along.

So after forking, switch to `starter`:
```bash
git clone https://github.com/<your-user>/yeet-to-prod
cd yeet-to-prod
git checkout starter        # the pre-CI starting point
```
> If you follow the guide on `main`, the workflows already exist — there's nothing to add.
> Always follow along on **`starter`**.

Two more things after forking:
1. **Enable Actions.** Forks have workflows **disabled** by default. Go to your fork →
   **Actions** tab → click **"I understand my workflows, go ahead and enable them"**.
   Nothing runs until you do this.
2. **Keep the repo name `yeet-to-prod`** so the `base: "/yeet-to-prod/"` in step A1 still
   matches. If you rename it, update `base` to `/<new-name>/`.

**Option B — Use your own project repo.** Skip forking. Apply the same steps to your repo;
just swap every `yeet-to-prod` value for yours (that's what **`← replace with yours`** means).

> Secrets and Variables are **not** copied when you fork — you'll add your own in the steps
> below. That's expected.

---

## 3. Part A — Frontend → GitHub Pages

### A1. Set the base path + test environment (`frontend/vite.config.js`)

GitHub Pages serves a project site under a subpath: `https://<user>.github.io/<repo>/`.
Vite must know this or the CSS/JS will 404. The same file also configures Vitest — the
component tests render React into a DOM, and Vitest's default environment is `node`, which
has no `document`. Edit `frontend/vite.config.js`:

```js
export default defineConfig(({ command }) => ({
  base: command === "build" ? "/yeet-to-prod/" : "/",   // ← replace yeet-to-prod with YOUR repo name
  plugins: [react()],
  test: {
    environment: "jsdom",   // component tests need a fake browser DOM
  },
}));
```

Why the `command ===` check: local `dev` stays at `/` (simple URLs); only the production
`build` gets the `/repo/` subpath.

Why `environment: "jsdom"`: `jsdom` is a browser-like DOM implemented in JavaScript. Without
it, `npm test` fails on `render(<App />)` with `ReferenceError: document is not defined`, and
since `pages.yml` runs `npm test` before building, **the whole deploy stops**. `jsdom` is
already in `devDependencies` — only this config line is missing.

Verify before you push:
```bash
cd frontend
npm install
npm test        # all tests must pass
```

### A2. Add the workflow file

Create `.github/workflows/pages.yml` (see the full file + explanation in
[section 6](#6-reference--the-workflow-scripts)). In short it does:
**lint + test on every push → build + deploy to Pages on `main`.**

### A3. Turn on Pages (GitHub settings)

Browser → your repo → **Settings → Pages** → under **Build and deployment**, set
**Source = "GitHub Actions"**. (Not "Deploy from a branch".)

### A4. Add the build-time variables

The frontend needs to know the backend URL at **build** time. Store it in GitHub (never
commit it):

Repo → **Settings → Secrets and variables → Actions → Variables tab → Repository
variables → New repository variable**:

| Name | Value | Note |
|---|---|---|
| `VITE_API_URL` | `https://yeet-to-prod.onrender.com` **← replace with yours** | your Render URL (from Part B) |
| `VITE_DEFAULT_COUNTRY` | `ID` | optional app setting |

> **Repository** variables (top section), not **Environment** variables — otherwise the
> build job can't read them. See Troubleshooting.
>
> Use **Variables**, not **Secrets** — this is a public URL, not a password.

### A5. Push and watch

```bash
git add .
git commit -m "add Pages deploy workflow"
git push
```
Repo → **Actions** tab → watch the run: `test → build → deploy`. When green, your site is
at:
```
https://qornanali.github.io/yeet-to-prod/     ← replace with yours
```

---

## 4. Part B — Backend → Render

### B1. Create the web service

<https://dashboard.render.com> → **New → Web Service** → connect your repo. Settings:

| Field | Value |
|---|---|
| **Root Directory** | `backend` |
| **Runtime** | Go |
| **Build Command** | `go build -o app .` |
| **Start Command** | `./app` |
| **Instance Type** | Free |

Create it. Render does a first deploy and gives you a URL like
`https://yeet-to-prod.onrender.com` **← replace with yours**.

### B2. Set the backend's environment variables

Render → your service → **Environment** → add:

| Key | Value | Why |
|---|---|---|
| `ALLOWED_ORIGIN` | `https://qornanali.github.io` **← replace with yours** | lets your Pages site call the API (CORS). **Host only — no path, no trailing slash.** |

### B3. Turn OFF Render auto-deploy

Render → Settings → **Auto-Deploy → No**. We want **GitHub Actions** to be the only thing
that deploys, so our tests act as the gate.

### B4. Get the Deploy Hook + store it in GitHub

Render → Settings → **Deploy Hook** → copy the URL (looks like
`https://api.render.com/deploy/srv-xxxx?key=yyyy`).

GitHub repo → **Settings → Secrets and variables → Actions → Secrets tab → New repository
secret**:

| Name | Value |
|---|---|
| `RENDER_DEPLOY_HOOK` | the hook URL you copied |

> Use a **Secret** here (not a Variable) — anyone with this URL can trigger a production
> deploy, so it must stay hidden.

### B5. Add the workflow file

Create `.github/workflows/backend.yml` (full file + explanation in
[section 6](#6-reference--the-workflow-scripts)). It does: **vet + test on every push →
deploy only when you click "Run workflow".**

### B6. Deploy it (manually)

```bash
git add .
git commit -m "add backend deploy workflow"
git push               # this runs tests only
```
Then to actually deploy: repo → **Actions** → "Backend CI + Deploy to Render" →
**Run workflow ▼** → branch `main` (optionally paste a commit SHA) → **Run workflow**.
Tests run, then Render redeploys.

Verify:
```bash
curl "https://yeet-to-prod.onrender.com/should-i-deploy?country=ID"     # ← replace with yours
```

---

## 5. Wire them together

The two sides must agree on two URLs — this is the #1 source of "why is it not working":

```
Frontend build   VITE_API_URL      = https://yeet-to-prod.onrender.com   (points AT backend)
Backend env      ALLOWED_ORIGIN    = https://qornanali.github.io         (lets frontend IN)
```

Both HTTPS. If you change either host, update both and re-run the frontend build (Actions →
Run workflow, or push a frontend change).

---

## 6. Reference — the workflow scripts

### `.github/workflows/pages.yml`

```yaml
name: Frontend CI + Deploy to Pages

on:
  push:
    paths:                       # only run when the frontend (or this file) changes
      - "frontend/**"
      - ".github/workflows/pages.yml"
  workflow_dispatch:             # also allow manual "Run workflow"

permissions:                     # what the deploy needs (least privilege)
  contents: read
  pages: write
  id-token: write

concurrency:                     # never run two Pages deploys at once
  group: pages
  cancel-in-progress: false

defaults:
  run:
    working-directory: frontend  # run all commands inside frontend/

jobs:
  test:                          # ── CI: runs on EVERY branch push ──
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
          cache-dependency-path: frontend/package-lock.json
      - run: npm ci               # install EXACT versions from the lockfile
      - run: npm run lint
      - run: npm test

  build:                         # ── produce the site. main only ──
    needs: test                  # only if tests passed
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 20, cache: npm, cache-dependency-path: frontend/package-lock.json }
      - run: npm ci
      - name: Build
        run: npm run build
        env:                     # inject the repo Variables at build time
          VITE_API_URL: ${{ vars.VITE_API_URL }}
          VITE_DEFAULT_COUNTRY: ${{ vars.VITE_DEFAULT_COUNTRY }}
      - uses: actions/configure-pages@v5
      - uses: actions/upload-pages-artifact@v3
        with:
          path: frontend/dist     # the built site

  deploy:                        # ── publish to Pages. main only ──
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
```

**Line-by-line:**

- **`on.push.paths`** — the workflow only triggers when files under `frontend/` (or the
  workflow itself) change. A backend-only commit won't rebuild the site.
- **`workflow_dispatch`** — adds a manual "Run workflow" button in the Actions tab.
- **`permissions`** — Pages deployment needs `pages: write` + `id-token: write`. We give
  the minimum, nothing more.
- **`concurrency: pages`** — if you push twice fast, deploys queue instead of racing.
- **`defaults.run.working-directory: frontend`** — so we don't write `cd frontend` in every
  step.
- **`jobs.test`** — runs on **every branch**. This is your quality gate: install deps, lint,
  test. `npm ci` installs the *exact* versions from `package-lock.json` (reproducible — the
  same on your machine and the robot's).
- **`if: github.ref == 'refs/heads/main'`** on `build`/`deploy` — build + publish happen
  **only on `main`**. Feature branches only get tested.
- **`needs: test` / `needs: build`** — the chain. `deploy` runs only if `build` ran, which
  runs only if `test` passed. **This is CI/CD in one keyword.**
- **`upload-pages-artifact` / `deploy-pages`** — the modern Pages flow: one job packages the
  built `dist/` as an artifact, the next publishes it.

### `.github/workflows/backend.yml`

```yaml
name: Backend CI + Deploy to Render

on:
  push:
    paths:
      - "backend/**"
      - ".github/workflows/backend.yml"
  workflow_dispatch:
    inputs:
      sha:                        # optionally deploy a specific commit
        description: "Commit SHA to test + deploy (blank = branch tip)"
        required: false

permissions:
  contents: read

concurrency:
  group: backend
  cancel-in-progress: false

defaults:
  run:
    working-directory: backend

jobs:
  test:                          # ── CI: every branch push ──
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{ inputs.sha || github.sha }}   # test the chosen commit
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: false            # stdlib-only app, nothing to cache
      - run: go vet ./...
      - run: go test ./...

  deploy:                        # ── MANUAL only, gated by test ──
    needs: test
    if: github.event_name == 'workflow_dispatch' && github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Trigger Render deploy
        working-directory: ${{ github.workspace }}   # no checkout here → don't use backend/
        run: curl -fsS "${{ secrets.RENDER_DEPLOY_HOOK }}&ref=${{ inputs.sha || github.sha }}"
```

**Line-by-line:**

- **`workflow_dispatch.inputs.sha`** — the "Run workflow" form shows a text box. Paste a
  commit SHA to deploy that exact commit; leave blank for the branch tip. Handy for
  **rollback** (redeploy a known-good older commit).
- **`test` checkout `ref: inputs.sha || github.sha`** — checks out the chosen commit, so the
  commit you **test** is the commit you **deploy** (no drift).
- **`deploy.if: github.event_name == 'workflow_dispatch'`** — deploy runs **only** via the
  manual button, never on a plain push. This makes shipping a human decision =
  **Continuous Delivery**.
- **`needs: test`** — even a manual deploy re-runs the tests first. Each workflow run is
  independent, so it can't reuse an earlier run's result — it re-verifies green, then ships.
- **`curl ...&ref=<sha>`** — calls the Render Deploy Hook. `-f` makes the step fail if
  Render returns an error. `&ref=` tells Render **which commit to build**.
- **`working-directory: ${{ github.workspace }}`** — the deploy job doesn't check out the
  repo (it only sends a URL), so the default `backend/` folder doesn't exist here; we point
  at the workspace root instead.

---

## 7. Troubleshooting (real errors)

**`npm test` fails: `ReferenceError: document is not defined` (the `test` job goes red, so
nothing deploys)**
```
FAIL  src/App.test.jsx > App > shows backend message + green bg when safe
ReferenceError: document is not defined
 ❯ Proxy.render node_modules/@testing-library/react/dist/pure.js:256:5
```
→ Vitest ran in its default `node` environment, which has no DOM. Add the `test` block from
[step A1](#a1-set-the-base-path--test-environment-frontendviteconfigjs) to
`frontend/vite.config.js`:
```js
test: {
  environment: "jsdom",
},
```
Pure-logic tests (`logic.test.js`) still pass — only the ones calling `render()` fail, which
is the tell.

**Frontend loads but is unstyled / blank, console shows 404 on `/assets/...`**
→ `base` in `vite.config.js` doesn't match your repo name. It must be `/<repo>/` for a
project site. Fix, push, rebuild.

**Browser console: "blocked by CORS policy"**
→ Backend `ALLOWED_ORIGIN` doesn't match. It must be the **origin only**:
`https://qornanali.github.io` — **no path** (`/yeet-to-prod/`) and **no trailing slash**
(`.io/`). The browser sends `Origin: https://qornanali.github.io`; your allowlist compares
it exactly. Verify:
```bash
curl -i -H "Origin: https://qornanali.github.io" \
  "https://yeet-to-prod.onrender.com/should-i-deploy?country=ID" | grep -i access-control
# want: access-control-allow-origin: https://qornanali.github.io
```

**Frontend calls `localhost:8080` in production**
→ `VITE_API_URL` wasn't set at build time. Most common cause: you added it as an
**Environment variable** instead of a **Repository variable**. The build job has no
`environment:`, so it can't read environment-scoped vars. Move it to **Repository
variables** and re-run the build.

**Deploy job fails: "No such file or directory ... working directory .../backend"**
→ A job with no `checkout` step inherited `working-directory: backend`, which doesn't exist
because the repo wasn't fetched. Add `working-directory: ${{ github.workspace }}` to that
step (as in `backend.yml` above).

**`npm ci` fails: "package-lock.json not found" / out of sync**
→ Commit your `frontend/package-lock.json`. `npm ci` needs it (that's the whole point —
exact, reproducible installs).

**Render backend is slow on first request**
→ Free instances sleep after ~15 min idle and cold-start (~30s). Normal. Hit the URL once to
wake it before a demo.

---

## 8. Glossary

- **CI (Continuous Integration)** — automatically build + test your code on every push. "Is
  it good?"
- **CD** — two flavors:
  - **Continuous Delivery** — every good build is ready to ship; a **human clicks** to
    release (our backend).
  - **Continuous Deployment** — every good build ships **automatically**, no human (our
    frontend).
- **Workflow** — a `.yml` file in `.github/workflows/` describing jobs the robots run.
- **Job / step** — a job is a group of steps on a fresh machine; steps run in order.
- **`needs:`** — makes one job wait for another. `deploy needs: test` = deploy only if tests
  passed.
- **`workflow_dispatch`** — a trigger that adds a manual "Run workflow" button.
- **Artifact** — a file bundle a job produces and hands to a later job (our built `dist/`).
- **Deploy Hook** — a secret URL; calling it tells Render to deploy. `?ref=<sha>` picks the
  commit.
- **Variable vs Secret** — both stored in GitHub settings, referenced in the yml (never
  committed). **Variables** are plain text (public config like a URL). **Secrets** are
  hidden (passwords, hook URLs).
- **CORS** — a browser rule: a page on one origin can only call another origin if that
  server allows it (`Access-Control-Allow-Origin`). Not real security — just browser hygiene.
- **Origin** — `scheme://host` only (e.g. `https://qornanali.github.io`). No path, no
  trailing slash.
- **Base path** — the subpath a site is served under. Pages project sites live at
  `/<repo>/`, so the frontend must build with that base.
