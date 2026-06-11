# Testing Checklist

## Server and Setup

- [ ] Server starts with `go run main.go`.
- [ ] Database connects.
- [ ] Migrations run.
- [ ] `/api/health` returns success.

## Multi-User Support

- [ ] Create user Joshua with consent.
- [ ] Add one goal for Joshua.
- [ ] Switch user (clear localStorage or use Switch button).
- [ ] Create user Bernard.
- [ ] Bernard does not see Joshua's goal.
- [ ] Add transaction for Bernard.
- [ ] Switch back to Joshua.
- [ ] Joshua does not see Bernard's transaction.
- [ ] Refresh page.
- [ ] Active user remains (localStorage).
- [ ] New browser/incognito starts from onboarding.

## Consent Flow

- [ ] Consent saves for a new user.
- [ ] Without consent, transactions, goals, affordability, Wizz are blocked.
- [ ] Consent status shows in the Privacy tab.

## Demo Data

- [ ] Load demo data replaces current user records (after confirmation).
- [ ] Refresh does not delete user data.
- [ ] User A loading demo does not affect User B.
- [ ] Load demo data button asks confirmation.
- [ ] Reset deletes only the current user's data.

## Transactions

- [ ] Transactions add and delete for the correct user.
- [ ] Dashboard updates after transaction changes.

## Goals

- [ ] Goals create and delete for the correct user.
- [ ] Goal progress and completion estimate display.

## Affordability

- [ ] High-risk result for SLE 7,000 laptop (low savings user).
- [ ] Low-risk result for affordable item (high savings user).
- [ ] Result uses the selected user's income, expenses, savings, goals.
- [ ] Zero surplus does not give a fake safe purchase date.

## Wizz Chat

- [ ] Wizz response uses the selected user's financial data.
- [ ] Wizz structured response (bullets, bold labels, spacing).
- [ ] Wizz fallback works without GROQ_API_KEY.
- [ ] Wizz works with GROQ_API_KEY.
- [ ] Wizz answer in Krio for "Wetin na osusu?"
- [ ] Safe formatting: **bold**, bullets, paragraphs render correctly.
- [ ] Typing indicator shows while waiting.
- [ ] Rate limiting after many requests.

## Dark Mode

- [ ] Toggle dark mode.
- [ ] Refresh page — dark mode persists.
- [ ] Cards, forms, tables, chat are readable in dark mode.
- [ ] Toggle back to light mode.
- [ ] Glassmorphism cards readable in both modes.

## JSON / CSV Import and Export

- [ ] JSON export includes profile, transactions, goals.
- [ ] JSON import restores profile, transactions, goals.
- [ ] Valid JSON backup imports cleanly.
- [ ] Invalid JSON shows clear error.
- [ ] JSON missing schemaVersion shows error.
- [ ] CSV export exports transactions only.
- [ ] CSV import imports transactions only.
- [ ] Valid CSV transaction file imports cleanly.
- [ ] CSV with wrong headers shows clear error.
- [ ] CSV with invalid data (negative amounts, bad dates) shows row-level errors.
- [ ] "CSV import is for transactions only" message visible in the UI.

## Sample Data

- [ ] `sample-data/demo-data.json` loads via JSON import.
- [ ] `sample-data/sample-transactions.csv` loads via CSV import.
- [ ] `sample-data/invalid-transactions.csv` shows clear validation errors.

## Reset

- [ ] Reset deletes only the active user's data.
- [ ] Other users' data remains after reset.
- [ ] Reset-all (development tools) deletes everything.

## Security

- [ ] No API key exposed in frontend files.
- [ ] No user-facing text mentions "Groq" (only technical docs).
- [ ] Wizz is the user-facing name everywhere.
- [ ] All SQL uses parameterized queries.
- [ ] Request body size is limited.
- [ ] Security headers present (X-Content-Type-Options, X-Frame-Options, Referrer-Policy).
- [ ] `.env` is in `.gitignore`.
- [ ] `SECURITY.md` describes how to report issues.

## Final Acceptance

1. Start server.
2. Open app.
3. Create user Joshua with consent.
4. Add one goal.
5. Refresh page — goal remains.
6. Switch to Bernard — no Joshua goal visible.
7. Toggle dark mode — works.
8. Refresh page — dark mode persists.
9. Add transactions for Bernard.
10. Dashboard updates.
11. Check affordability for Laptop SLE 7000 — uses Bernard's data.
12. Ask Wizz "Can I buy this laptop?" — structured response with data context.
13. Ask Wizz "Wetin na osusu?" — Krio/English answer.
14. Export JSON.
15. Export CSV.
16. Import valid JSON.
17. Import invalid CSV — clear error shown.
18. Reset Bernard only.
19. Switch back to Joshua — data still available.
20. No user-facing screen exposes GROQ_API_KEY or raw provider details.
