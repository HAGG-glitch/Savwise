# Architecture

## Current MVP

The current SavWise AI MVP uses HTML, Tailwind CSS, vanilla JavaScript, a Go backend, PostgreSQL storage, JSON/CSV export, and Groq for educational AI coaching.

The frontend is served by the Go application from the `web/` folder. JavaScript calls JSON API endpoints under `/api`. The backend performs validation, database operations, financial calculations, CSV generation, and Groq requests.

## Future production phase

Future work may include authentication, secure sessions, role-based access, encrypted fields, audit logging, Orange Money and Afrimoney provider-approved APIs, USSD/SMS integrations, and reviewed AI safety controls.
