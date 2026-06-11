# SavWise AI

SavWise AI is an open-source full-stack MVP prototype for DLAW207 I.T Law and IPR Legal Issues. It supports **SDG 1: No Poverty** by helping Sierra Leonean users track spending, build savings goals, check purchase affordability, and receive educational AI financial coaching.

## Technology stack

- Frontend: HTML, Tailwind CSS CDN, vanilla JavaScript
- Backend: Go `net/http`
- Database: PostgreSQL
- AI coach: Wizz (powered by Groq Chat Completions API through the Go backend)
- Data export: JSON and CSV

## Academic honesty boundary

This is a university MVP prototype. It uses manually entered or simulated financial data. It does not connect to real Orange Money or Afrimoney accounts, transfer money, perform KYC, issue loans, or provide regulated investment, legal, banking, or professional financial advice.

## Security Notes

- **Never commit `.env`** to version control. It contains your `GROQ_API_KEY`.
- If your `GROQ_API_KEY` is exposed, rotate it immediately in your Groq console.
- This is a university MVP. It does not process real payments or connect to real mobile-money accounts.
- Do not enter PINs, OTPs, passwords, national ID numbers, or real KYC data.
- All SQL queries use parameterized inputs. User input is validated before processing.

## Setup

1. Install Go and PostgreSQL.
2. Create a database:

```bash
createdb savwise_ai
```

3. Copy `.env.example` to `.env` and edit the values:

```bash
cp .env.example .env
```

4. Install Go dependencies:

```bash
go mod tidy
```

5. Run the app from the project root:

```bash
go run main.go
```

6. Open:

```text
http://localhost:8080
```

Migrations run automatically from the `migrations/` folder on startup.

## Current MVP Limitations

- Multi-user support is prototype-level using localStorage. Not production-ready authentication.
- No real mobile-money API integration (requires provider approval).
- CSV import supports transactions only. Use JSON for full backup/restore.
- Dark mode uses custom CSS variables (Tailwind class strategy is limited with CDN).
- No HTTPS in development. Production deployment must add TLS.
- No automated test suite yet.

## Wizz (AI Coach) setup

Wizz is the financial coach name shown to users. It uses Groq Chat Completions API through the Go backend.

Set `GROQ_API_KEY` in `.env`. If the key is missing or the service is unavailable, Wizz returns a local rule-based fallback response.

## Main API routes

- `GET /api/health`
- `GET /api/profile`
- `POST /api/profile`
- `POST /api/consent`
- `GET /api/transactions`
- `POST /api/transactions`
- `DELETE /api/transactions/{id}`
- `GET /api/goals`
- `POST /api/goals`
- `PUT /api/goals/{id}`
- `DELETE /api/goals/{id}`
- `GET /api/dashboard`
- `POST /api/affordability`
- `POST /api/coach`
- `GET /api/export/json`
- `GET /api/export/csv`
- `POST /api/import/json`
- `POST /api/load-demo`
- `DELETE /api/reset`

## Team members

- Joshua Yoki
- Bernard
- Moses Moore

## Licence

MIT License. See `LICENSE`.

## If Go proxy is blocked

Try direct module download:

```bash
go env -w GOPROXY=direct
go mod tidy
```

If direct download is also blocked, connect to a network that allows GitHub/module downloads, run `go mod tidy` once, then keep the generated `go.sum` in the project.

## Suggested Improvements After MVP Review

1. Add real authentication for multiple users.
2. Add role-based admin access.
3. Add stronger server-side validation and audit logging.
4. Add secure mobile-money API integration only after provider approval.
5. Add PostgreSQL row ownership checks for multi-user mode.
6. Add better chart visualisations for spending trends.
7. Add CSV import confirmation before saving imported data.
8. Add language toggle for English and Krio.
9. Add reviewed financial-literacy lessons.
10. Add AI safety testing to check for harmful or misleading advice.
11. Add mobile-first improvements for market traders.
12. Add offline mode later if possible.
13. Add automated tests for affordability and dashboard calculations.
14. Add deployment instructions for Render, Railway, or Fly.io.
15. Add a future USSD-lite flow for feature-phone users.
