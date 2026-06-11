# Security Policy

## Current MVP Status

SavWise AI is a university MVP prototype. It is **not** a production financial system. It does not process real payments, connect to real mobile-money accounts, or store sensitive personal financial data.

## Reporting a Vulnerability

If you discover a security issue in this prototype, please report it by emailing the development team. Do not submit real credentials, API keys, or personal data in any report.

## Security Practices

- Keep `.env` private. Never commit it to version control.
- If the `GROQ_API_KEY` is exposed, rotate it immediately in the Groq console.
- This prototype does not process real mobile-money transactions.
- All SQL queries use parameterized inputs.
- User input is validated before processing.
- The AI coach (Wizz) does not ask for PINs, OTPs, passwords, national ID numbers, or real KYC data.

## Future Production Security Requirements

For any production deployment, the following must be addressed:

- Authentication and role-based access control
- Encryption at rest and in transit
- HTTPS deployment with valid certificates
- Audit logging for all privacy-sensitive actions
- Regular backups
- Penetration testing
- Provider-approved API integration
- Privacy impact assessment
- Formal threat modelling
