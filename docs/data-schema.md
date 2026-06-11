# Data Schema

## User/Profile

Fields: `id`, `fullName`, `email`, `preferredLanguage`, `monthlyIncome`, `currentSavings`, `emergencyTarget`, `consentAccepted`.

## Transaction

Fields: `id`, `userId`, `description`, `amount`, `type`, `category`, `date`, `createdAt`.

Allowed `type`: `income`, `expense`.

CSV headers: `id,description,amount,type,category,date,createdAt`.

## Goal

Fields: `id`, `userId`, `name`, `targetAmount`, `currentAmount`, `monthlyContribution`, `progressPercent`, `remainingAmount`, `estimatedCompletion`, `status`.

## JSON export

```json
{
  "schemaVersion": "1.0",
  "exportedAt": "ISO_DATE",
  "user": {},
  "transactions": [],
  "goals": []
}
```
