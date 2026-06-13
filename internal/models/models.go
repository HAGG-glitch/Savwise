package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type User struct {
	ID                string  `json:"id"`
	FullName          string  `json:"fullName"`
	Email             string  `json:"email"`
	PreferredLanguage string  `json:"preferredLanguage"`
	MonthlyIncome     float64 `json:"monthlyIncome"`
	CurrentSavings    float64 `json:"currentSavings"`
	EmergencyTarget   float64 `json:"emergencyTarget"`
	ConsentAccepted   bool    `json:"consentAccepted"`
}

type CreateUserRequest struct {
	FullName          string  `json:"fullName"`
	Email             string  `json:"email"`
	PreferredLanguage string  `json:"preferredLanguage"`
	MonthlyIncome     float64 `json:"monthlyIncome"`
	CurrentSavings    float64 `json:"currentSavings"`
	EmergencyTarget   float64 `json:"emergencyTarget"`
	ConsentAccepted   bool    `json:"consentAccepted"`
}

type ProfileRequest struct {
	UserID            string  `json:"user_id"`
	FullName          string  `json:"fullName"`
	Email             string  `json:"email"`
	PreferredLanguage string  `json:"preferredLanguage"`
	MonthlyIncome     float64 `json:"monthlyIncome"`
	CurrentSavings    float64 `json:"currentSavings"`
	EmergencyTarget   float64 `json:"emergencyTarget"`
	ConsentAccepted   bool    `json:"consentAccepted"`
}

type Transaction struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	CreatedAt   string  `json:"createdAt"`
}

type TransactionRequest struct {
	UserID      string  `json:"user_id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}

type Goal struct {
	ID                  string  `json:"id"`
	UserID              string  `json:"userId"`
	Name                string  `json:"name"`
	TargetAmount        float64 `json:"targetAmount"`
	CurrentAmount       float64 `json:"currentAmount"`
	MonthlyContribution float64 `json:"monthlyContribution"`
	CreatedAt           string  `json:"createdAt"`
	ProgressPercent     float64 `json:"progressPercent,omitempty"`
	RemainingAmount     float64 `json:"remainingAmount,omitempty"`
	EstimatedCompletion string  `json:"estimatedCompletion,omitempty"`
	Status              string  `json:"status,omitempty"`
}

type GoalRequest struct {
	UserID              string  `json:"user_id"`
	Name                string  `json:"name"`
	TargetAmount        float64 `json:"targetAmount"`
	CurrentAmount       float64 `json:"currentAmount"`
	MonthlyContribution float64 `json:"monthlyContribution"`
}

type ScoreBreakdown struct {
	SavingsHabit  int `json:"savingsHabit"`
	BudgetControl int `json:"budgetControl"`
	EmergencyFund int `json:"emergencyFund"`
	GoalProgress  int `json:"goalProgress"`
	Total         int `json:"total"`
}

type Alert struct {
	Severity          string `json:"severity"`
	Title             string `json:"title"`
	Explanation       string `json:"explanation"`
	RecommendedAction string `json:"recommendedAction"`
}

type CategoryTotal struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Percent  float64 `json:"percent"`
}

type Dashboard struct {
	User                   User            `json:"user"`
	TotalIncome            float64         `json:"totalIncome"`
	TotalExpenses          float64         `json:"totalExpenses"`
	MonthlySurplus         float64         `json:"monthlySurplus"`
	SavingsRate            float64         `json:"savingsRate"`
	EmergencyCoverageDays  float64         `json:"emergencyCoverageDays"`
	BiggestExpenseCategory string          `json:"biggestExpenseCategory"`
	SpendingBreakdown      []CategoryTotal `json:"spendingBreakdown"`
	Score                  ScoreBreakdown  `json:"score"`
	Alerts                 []Alert         `json:"alerts"`
	Goals                  []Goal          `json:"goals"`
	SavingsOpportunity     string          `json:"savingsOpportunity"`
	WeeklyReport           string          `json:"weeklyReport"`
}

type AffordabilityRequest struct {
	UserID     string  `json:"user_id"`
	ItemName   string  `json:"itemName"`
	ItemPrice  float64 `json:"itemPrice"`
	TargetDate string  `json:"targetDate,omitempty"`
}

type AffordabilityResult struct {
	ID                    string   `json:"id,omitempty"`
	ItemName              string   `json:"itemName"`
	ItemPrice             float64  `json:"itemPrice"`
	TargetDate            string   `json:"targetDate"`
	CalculatedAt          string   `json:"calculatedAt"`
	ExpensePeriod         string   `json:"expensePeriod"`
	MonthlyIncome         float64  `json:"monthlyIncome"`
	MonthlyExpenses       float64  `json:"monthlyExpenses"`
	MonthlySurplus        float64  `json:"monthlySurplus"`
	CurrentSavings        float64  `json:"currentSavings"`
	EmergencyTarget       float64  `json:"emergencyTarget"`
	FundingGap            float64  `json:"fundingGap"`
	MonthsUntilTarget     int      `json:"monthsUntilTarget"`
	RequiredMonthlySaving float64  `json:"requiredMonthlySaving"`
	ActiveGoalCommitments float64  `json:"activeGoalCommitments"`
	AvailableAfterGoals   float64  `json:"availableAfterGoals"`
	RiskLevel             string   `json:"riskLevel"`
	Reasons               []string `json:"reasons"`
	GoalImpact            string   `json:"goalImpact"`
	Recommendation        string   `json:"recommendation"`
	Explanation           string   `json:"explanation"`
}

type CoachRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type CoachResponse struct {
	Response string `json:"response"`
	Source   string `json:"source"`
	Model    string `json:"model"`
}

type ExportPackage struct {
	SchemaVersion string        `json:"schemaVersion"`
	ExportedAt    string        `json:"exportedAt"`
	User          User          `json:"user"`
	Transactions  []Transaction `json:"transactions"`
	Goals         []Goal        `json:"goals"`
}
