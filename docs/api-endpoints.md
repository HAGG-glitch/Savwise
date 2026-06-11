# API Endpoints

All JSON endpoints return:

```json
{ "success": true, "message": "...", "data": {} }
```

or:

```json
{ "success": false, "message": "...", "error": "code" }
```

Routes:

- `GET /api/health` — check server/database status.
- `GET /api/profile` — load demo profile.
- `POST /api/profile` — save profile and consent.
- `POST /api/consent` — record consent.
- `GET /api/transactions` — list transactions.
- `POST /api/transactions` — add transaction.
- `DELETE /api/transactions/{id}` — delete transaction.
- `GET /api/goals` — list goals.
- `POST /api/goals` — create goal.
- `PUT /api/goals/{id}` — update goal.
- `DELETE /api/goals/{id}` — delete goal.
- `GET /api/dashboard` — calculate dashboard.
- `POST /api/affordability` — calculate purchase risk.
- `POST /api/coach` — ask Groq or fallback coach.
- `GET /api/export/json` — export user data.
- `GET /api/export/csv` — export transactions.
- `POST /api/import/json` — import JSON backup.
- `POST /api/load-demo` — load demo data.
- `DELETE /api/reset` — delete demo data.
