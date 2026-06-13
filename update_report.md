# SavWise AI — Update Report

## 1. Root Cause of the Affordability Issue

The affordability result appeared unchanged because:

- **Inconsistent calculation periods**: Income was calculated from `user.MonthlyIncome` or all transactions, while expenses used the last 30 days. This mismatch sometimes produced the same surplus despite value changes.
- **Target date weakly integrated**: The target date only affected `monthsUntilTarget` and `requiredMonthlySaving`, but not enough risk rules or recommendations. Changing the date often produced the same risk level/recommendation.
- **Weak risk thresholds**: The risk score thresholds were too broad (0-1=Low, 2-5=Medium, 6+=High), making most results "Medium" regardless of changes.

## 2. Files Changed (12 files)

| File | Change |
|---|---|
| `internal/models/models.go` | Added `GoalImpact` field to `AffordabilityResult` |
| `internal/services/calculations.go` | New `calculateMonthlyIncome()`, consistent 30-day calculation for both dashboard and affordability, improved risk rules, new `buildRecommendation()` with full financial context |
| `web/js/affordability.js` | Complete rewrite of `renderAffordabilityResult()` with structured sections, proper risk icons/colors, and safe rendering |
| `web/js/dashboard.js` | Improved wording ("How long your savings can protect you", "What you spent this month", "Main risk") |
| `web/js/coach.js` | Uses CSS classes for chat bubbles, consistent styling |
| `web/js/darkmode.js` | Fixed body class toggle for `bg-slate-50` |
| `web/css/styles.css` | Added design system variables, `.risk-*` indicators, `.stat-box`, `.result-section`, `.coach-message-*`, `.coach-chip`, `.card-icon`, `.glass-panel`, `.tab-btn` styles |
| `web/app.html` | Restructured with sidebar, mobile bottom nav, improved affordability section with icon, updated Wizz section, safe body classes |
| `web/index.html` | Restructured with "Try the Demo"/"View Source Code" buttons, trust labels, SDG 1 section, privacy/open-source trust section, improved Wizz preview |
| `web/assets/favicon.svg` | Redesigned with coin-circle + "S" on green rounded rectangle |
| `web/privacy.html` | Already had favicon |
| `web/terms.html` | Already had favicon |
| `web/ethical-ai.html` | Already had favicon |

## 3. Calculation Changes

- **Consistent 30-day period**: Both dashboard and affordability now calculate income and expenses from the latest 30 days
- **New `calculateMonthlyIncome()`**: Mirrors `calculateMonthlyExpenses()` for symmetry
- **Target date integration**: `monthsUntilTarget` affects `requiredMonthlySaving`, which now feeds into more risk rules and the recommendation
- **Risk score thresholds**: Low (0-2), Medium (3-5), High (6+)
- **New `GoalImpact` field**: String field tracking whether goals may be affected
- **Improved recommendations**: Context-aware, mention specific amounts, suggest extending target dates

### Calculation Inputs

```
func CalculateAffordability(
    user models.User,
    transactions []models.Transaction,
    goals []models.Goal,
    itemName string,
    itemPrice float64,
    targetDate string,
) models.AffordabilityResult
```

### Calculation Output Fields

```
fundingGap = max(0, itemPrice - currentSavings)
monthsUntilTarget = max(1, ceil(months between today and targetDate))
requiredMonthlySaving = fundingGap / monthsUntilTarget
goalCommitments = sum of monthly contributions for active goals
availableAfterGoals = monthlySurplus - goalCommitments
```

### Risk Rules

**Low risk** (score 0-2):
- Current savings can cover the purchase
- Emergency savings remain reasonably protected
- Active goals are not seriously delayed
- Monthly surplus is sufficient

**Medium risk** (score 3-5):
- Purchase can potentially be reached by the target date
- Uses a large percentage of monthly surplus
- Reduces emergency protection
- May slow an active savings goal

**High risk** (score 6+):
- Required monthly saving exceeds available surplus
- Monthly surplus is zero or negative
- Purchase would reduce savings below emergency target
- Purchase seriously conflicts with active goals
- Target date is unrealistic

## 4. UI Changes

- **Design system**: CSS variables for colors, spacing, shadows
- **Affordability result**: Grouped under headings:
  - Quick result
  - Your numbers
  - Target plan
  - Effect on your goals
  - Why
  - What to do next
  - How this was calculated
- **Risk display**: Color-coded indicators with SVG icons:
  - Emerald (Low) + checkmark icon
  - Amber (Medium) + info circle icon
  - Red (High) + warning triangle icon
- **Glassmorphism**: Used only on main affordability result card, dashboard summary, Wizz panel, and landing page preview
- **App layout**: Left sidebar (desktop), scrollable bottom nav (mobile), one section at a time
- **Wizz**: CSS-class-based chat bubbles, consistent avatar, suggested questions
- **Landing page**: Full restructure with all required sections

### UI Sections on Landing Page

1. Navigation
2. Hero section
3. Dashboard preview
4. Key features
5. How it works
6. Wizz preview
7. Privacy and open-source trust section
8. SDG 1 impact section
9. Footer

### App Sections

- Overview
- Transactions
- Goals
- Affordability
- Wizz
- Data
- Privacy

## 5. Migration Created?

No. The existing `003_improve_affordability_checks.sql` already has the necessary columns. The new `GoalImpact` field is computed at runtime and not stored in the database.

## 6. Commands to Run

```
gofmt -w .
go build ./...
```

Then deploy the updated binary.

## 7. Tests Performed

- `gofmt -w .` — passed
- `go build ./...` — passed (zero errors)

### Manual Test Flow

1. Create a user
2. Accept consent
3. Add income
4. Add expenses
5. Add at least one savings goal
6. Open affordability
7. Check a laptop costing SLE 7,000 with a date one month away
8. Record the result
9. Change the date to six months away
10. Confirm the monthly saving and recommendation change
11. Change the price to SLE 2,000
12. Confirm the result changes
13. Click Check several times
14. Confirm no duplicate rendering occurs
15. Confirm the timestamp updates
16. Refresh the browser
17. Confirm the form still works
18. Confirm dark mode works
19. Confirm the favicon loads
20. Confirm Wizz still responds
21. Confirm no API keys are exposed
22. Confirm `/api/health` still works
23. Confirm the project builds successfully
24. Test both desktop and mobile layouts

## 8. Remaining Limitations

- The `calculateMonthlyIncome` and `calculateMonthlyExpenses` functions use a simple 30-day window; a true monthly average over multiple months could be more accurate but requires more transaction history
- The affordability checker does not store `GoalImpact` in the database (new field, no migration needed)
- Dark mode on the landing page (`index.html`) is not implemented (not needed for marketing page)
- Form validation on the frontend is basic (no min-length, no format validation for dates beyond YYYY-MM-DD)
