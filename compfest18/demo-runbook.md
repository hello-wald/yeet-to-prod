# CI/CD Live Demo — Speaker Runbook

Pointer for running the live deploy demo. Follows **`deploy-guide.md`** step-by-step, on a
branch that starts *before* any CI existed — so adding the workflows is real, live work.

- **Repo:** `yeet-to-prod`
- **Demo branch:** `demo/deploy-guide` (branched from commit `3469423`, the last pre-CI
  commit — full app, no `.github/workflows/`, no vite `base`).
- **Reference:** `deploy-guide.md` (share this with students too).

---

## Pre-flight (do BEFORE the session, once)

Repo-level settings persist across branches, so set these on the repo now:

- [ ] **Pages** → Settings → Pages → Source = **GitHub Actions**
- [ ] **Repo Variables** → `VITE_API_URL`, `VITE_DEFAULT_COUNTRY` (Repository, not Environment)
- [ ] **Render** service exists · Root Dir `backend` · Auto-Deploy **OFF** · `ALLOWED_ORIGIN`
      includes `https://qornanali.github.io`
- [ ] **Repo Secret** → `RENDER_DEPLOY_HOOK`
- [ ] **Pre-warm Render** (hit the URL) ~5 min before — free instance cold-starts ~30s
- [ ] `deploy-guide.md` open in a tab; terminal on branch `demo/deploy-guide`
- [ ] Confirm you're on the branch: `git branch --show-current` → `demo/deploy-guide`

> Leave the demo branch **pre-CI** (don't pre-add the workflows) — adding them live is the demo.

---

## The flow (maps to deploy-guide sections)

### 0. Show the starting point (30 sec)
- `git branch --show-current` → `demo/deploy-guide`
- `ls .github/workflows` → nothing. *"No robot yet. Every push right now is unchecked."*

### 1. Frontend CI — add + watch it run  (guide §3 A1–A2, §6)
1. Edit `frontend/vite.config.js` → add the `base` (guide **A1**).
2. Add `.github/workflows/pages.yml` (guide **§6** — walk the key lines: `on.push`,
   `needs:`, the `if: main` gate).
3. Push the **branch**:
   ```bash
   git add . && git commit -m "add Pages CI/CD workflow" && git push -u origin demo/deploy-guide
   ```
4. Actions tab → **only `test` runs** (branch ≠ main → build/deploy skipped).
   → **Teaching point:** *"CI runs on every push. Deploy is gated to main."*

### 2. (Optional) Break it → red → fix → green  (the peak)
- Logic break: in `frontend/src/logic.js` (or backend `decide.go`) invert something → push →
  `test` goes **red** → *"the robot caught it before users did."* → revert → green.
- Or dependency break: delete `frontend/package-lock.json` → `npm ci` fails → commit lockfile → green.

### 3. Ship the frontend — merge to main  (guide §3 A3–A5)
- Merge/PR `demo/deploy-guide` → `main` (or push to main).
- Now `test → build → deploy` all run → site live at
  `https://qornanali.github.io/yeet-to-prod/`.
  → **Teaching point:** *"Push to main = auto-ship. That's Continuous Deployment."*

### 4. Backend — add + manual deploy  (guide §4, §6)
1. Add `.github/workflows/backend.yml` (guide **§6** — highlight `workflow_dispatch` +
   `deploy needs: test`).
2. Push → **only `test` runs** (no deploy on push).
3. Actions → **Run workflow** → branch `main` → Render redeploys.
   → **Teaching point:** *"Deploy waits for a human click. That's Continuous Delivery."*

### 5. Payoff
- Open the live Pages site → it calls the live Render backend → green/red verdict.
  *"Push to ship. Robots test + deploy. That's CI/CD."*

---

## Gotchas live (pre-answers)
- **Push on the demo branch doesn't deploy** — expected; deploy `if: main`. Merge to main to ship.
- **Frontend shows red/error** — backend cold-starting or `VITE_API_URL`/CORS mismatch. See
  `deploy-guide.md` §7 Troubleshooting.
- **Fork demo (if students follow):** Actions disabled by default on forks → enable in Actions tab.

---

## Reset between rehearsals
```bash
git checkout main
git branch -D demo/deploy-guide
git checkout -b demo/deploy-guide 3469423
git push origin :demo/deploy-guide 2>/dev/null || true   # delete remote branch if pushed
```
Back to a clean pre-CI branch.
