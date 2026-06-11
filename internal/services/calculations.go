package services

import (
	"math"
	"sort"
	"strings"
	"time"

	"savwise-ai/internal/models"
)

func EnhanceGoals(goals []models.Goal) []models.Goal {
	out := make([]models.Goal, 0, len(goals))
	for _, g := range goals {
		g.ProgressPercent = 0
		if g.TargetAmount > 0 {
			g.ProgressPercent = math.Min(100, (g.CurrentAmount/g.TargetAmount)*100)
		}
		g.RemainingAmount = math.Max(0, g.TargetAmount-g.CurrentAmount)
		switch {
		case g.ProgressPercent >= 100:
			g.Status = "Completed"
			g.EstimatedCompletion = "Completed"
		case g.CurrentAmount <= 0:
			g.Status = "Not started"
		case g.MonthlyContribution <= 0:
			g.Status = "Needs contribution"
			g.EstimatedCompletion = "No estimate"
		default:
			g.Status = "In progress"
			months := int(math.Ceil(g.RemainingAmount / g.MonthlyContribution))
			if months <= 0 {
				g.EstimatedCompletion = "Completed"
			} else {
				g.EstimatedCompletion = time.Now().AddDate(0, months, 0).Format("January 2006")
			}
		}
		out = append(out, g)
	}
	return out
}

func CalculateDashboard(user models.User, transactions []models.Transaction, goals []models.Goal) models.Dashboard {
	var totalIncome, totalExpenses float64
	categoryTotals := map[string]float64{}
	for _, tx := range transactions {
		if tx.Type == "income" {
			totalIncome += tx.Amount
		} else if tx.Type == "expense" {
			totalExpenses += tx.Amount
			categoryTotals[tx.Category] += tx.Amount
		}
	}
	if totalIncome == 0 && user.MonthlyIncome > 0 {
		totalIncome = user.MonthlyIncome
	}
	surplus := totalIncome - totalExpenses
	savingsRate := 0.0
	if totalIncome > 0 {
		savingsRate = (surplus / totalIncome) * 100
	}

	dailyExpenses := totalExpenses / 30
	coverage := 0.0
	if dailyExpenses > 0 {
		coverage = user.CurrentSavings / dailyExpenses
	} else if user.CurrentSavings > 0 {
		coverage = 90
	}

	breakdown := make([]models.CategoryTotal, 0, len(categoryTotals))
	biggest := "None"
	biggestValue := 0.0
	for cat, amt := range categoryTotals {
		pct := 0.0
		if totalExpenses > 0 {
			pct = (amt / totalExpenses) * 100
		}
		breakdown = append(breakdown, models.CategoryTotal{Category: cat, Amount: amt, Percent: pct})
		if amt > biggestValue {
			biggest = cat
			biggestValue = amt
		}
	}
	sort.Slice(breakdown, func(i, j int) bool { return breakdown[i].Amount > breakdown[j].Amount })

	enhancedGoals := EnhanceGoals(goals)
	score := calculateScore(totalIncome, totalExpenses, savingsRate, coverage, enhancedGoals)
	alerts := calculateAlerts(totalIncome, totalExpenses, savingsRate, coverage, categoryTotals, enhancedGoals)

	return models.Dashboard{
		User: user, TotalIncome: round(totalIncome), TotalExpenses: round(totalExpenses), MonthlySurplus: round(surplus),
		SavingsRate: round(savingsRate), EmergencyCoverageDays: round(coverage), BiggestExpenseCategory: biggest,
		SpendingBreakdown: breakdown, Score: score, Alerts: alerts, Goals: enhancedGoals,
		SavingsOpportunity: savingsOpportunity(categoryTotals), WeeklyReport: weeklyReport(totalIncome, totalExpenses, surplus, biggest),
	}
}

func calculateScore(income, expenses, savingsRate, coverage float64, goals []models.Goal) models.ScoreBreakdown {
	savingsHabit := 0
	switch {
	case savingsRate >= 20:
		savingsHabit = 30
	case savingsRate >= 10:
		savingsHabit = 20
	case savingsRate > 0:
		savingsHabit = 10
	}
	budgetControl := 0
	expensePct := 999.0
	if income > 0 {
		expensePct = (expenses / income) * 100
	}
	switch {
	case expensePct <= 70:
		budgetControl = 30
	case expensePct <= 90:
		budgetControl = 20
	case expensePct <= 100:
		budgetControl = 10
	}
	emergencyFund := 0
	switch {
	case coverage >= 90:
		emergencyFund = 20
	case coverage >= 30:
		emergencyFund = 12
	case coverage >= 7:
		emergencyFund = 6
	}
	goalProgress := 0
	if len(goals) > 0 {
		sum := 0.0
		for _, g := range goals {
			sum += math.Min(100, g.ProgressPercent)
		}
		avg := sum / float64(len(goals))
		goalProgress = int(math.Round((avg / 100) * 20))
	}
	total := savingsHabit + budgetControl + emergencyFund + goalProgress
	return models.ScoreBreakdown{SavingsHabit: savingsHabit, BudgetControl: budgetControl, EmergencyFund: emergencyFund, GoalProgress: goalProgress, Total: total}
}

func calculateAlerts(income, expenses, savingsRate, coverage float64, categoryTotals map[string]float64, goals []models.Goal) []models.Alert {
	var alerts []models.Alert
	if expenses > income && income > 0 {
		alerts = append(alerts, models.Alert{Severity: "High", Title: "Expenses exceed income", Explanation: "You are spending more than you earn in this period.", RecommendedAction: "Pause non-essential spending and review your budget immediately."})
	}
	if coverage < 7 {
		alerts = append(alerts, models.Alert{Severity: "High", Title: "Emergency fund is very low", Explanation: "Your current savings may not cover one week of expenses.", RecommendedAction: "Prioritise emergency savings before large purchases."})
	}
	if income > 0 && savingsRate < 10 {
		alerts = append(alerts, models.Alert{Severity: "Medium", Title: "Savings rate below target", Explanation: "Your savings rate is below the prototype target of 10%.", RecommendedAction: "Try redirecting small discretionary expenses to savings."})
	}
	for cat, amount := range categoryTotals {
		if income > 0 && (amount/income)*100 > 35 && !strings.EqualFold(cat, "Healthcare") {
			alerts = append(alerts, models.Alert{Severity: "Medium", Title: cat + " spending is high", Explanation: "This category uses more than 35% of your income.", RecommendedAction: "Review this category and set a spending limit."})
		}
	}
	if income > 0 && savingsRate >= 20 {
		alerts = append(alerts, models.Alert{Severity: "Positive", Title: "Savings contributions are on track", Explanation: "You are saving at least 20% of income in this period.", RecommendedAction: "Maintain this pattern and protect your emergency fund."})
	}
	for _, g := range goals {
		if g.ProgressPercent >= 100 {
			alerts = append(alerts, models.Alert{Severity: "Positive", Title: g.Name + " completed", Explanation: "This savings goal has reached 100%.", RecommendedAction: "Consider starting or strengthening an emergency fund."})
		} else if g.ProgressPercent >= 50 {
			alerts = append(alerts, models.Alert{Severity: "Positive", Title: g.Name + " passed 50%", Explanation: "You are more than halfway to this goal.", RecommendedAction: "Continue the monthly contribution until completion."})
		}
	}
	if len(alerts) == 0 {
		alerts = append(alerts, models.Alert{Severity: "Positive", Title: "No major risk detected", Explanation: "Current prototype rules did not detect urgent spending risks.", RecommendedAction: "Continue tracking income, expenses, and goals."})
	}
	return alerts
}

func savingsOpportunity(categoryTotals map[string]float64) string {
	if v := categoryTotals["Airtime/Data"]; v > 0 {
		return "Reducing Airtime/Data spending by 10% could save about SLE " + money(v*0.10) + " per month."
	}
	if v := categoryTotals["Entertainment"]; v > 0 {
		return "Reducing Entertainment spending by 10% could save about SLE " + money(v*0.10) + " per month."
	}
	return "Track more expense categories to discover savings opportunities."
}

func weeklyReport(income, expenses, surplus float64, biggest string) string {
	return "Income: SLE " + money(income) + "; Spent: SLE " + money(expenses) + "; Saved/surplus: SLE " + money(surplus) + "; Biggest expense: " + biggest + "."
}

func CalculateAffordability(user models.User, transactions []models.Transaction, itemName string, itemPrice float64) models.AffordabilityResult {
	var monthlyExpenses float64
	for _, tx := range transactions {
		if tx.Type == "expense" {
			monthlyExpenses += tx.Amount
		}
	}
	income := user.MonthlyIncome
	if income <= 0 {
		for _, tx := range transactions {
			if tx.Type == "income" {
				income += tx.Amount
			}
		}
	}
	surplus := income - monthlyExpenses
	remainingSavings := user.CurrentSavings - itemPrice
	reasons := []string{}
	riskScore := 0
	impactGoals := false

	if itemPrice > user.CurrentSavings {
		riskScore += 3
		reasons = append(reasons, "This item costs more than your current savings.")
	} else {
		reasons = append(reasons, "Your current savings can cover the price.")
	}
	if monthlyExpenses > 0 && remainingSavings < monthlyExpenses {
		riskScore += 3
		reasons = append(reasons, "After buying, remaining savings would be below one month of expenses.")
	} else if monthlyExpenses > 0 && remainingSavings < monthlyExpenses*3 {
		riskScore += 1
		reasons = append(reasons, "After buying, remaining savings would be below three months of expenses.")
	}
	if income > 0 && itemPrice > income*2 {
		riskScore += 3
		reasons = append(reasons, "This item costs more than twice your monthly income.")
	} else if income > 0 && itemPrice > income*0.5 {
		riskScore += 1
		reasons = append(reasons, "This item costs more than half of your monthly income.")
	}
	if surplus <= 0 {
		riskScore += 3
		reasons = append(reasons, "Your budget does not produce a positive monthly surplus.")
	}
	if user.CurrentSavings > 0 && itemPrice > user.CurrentSavings*0.5 {
		riskScore += 1
		reasons = append(reasons, "This purchase would use more than half of your savings.")
	}
	if user.EmergencyTarget > 0 && remainingSavings < user.EmergencyTarget {
		riskScore += 2
		impactGoals = true
		reasons = append(reasons, "This purchase would leave your emergency fund below target.")
	}

	risk := "Low"
	var wait *int

	if surplus > 0 && itemPrice > user.CurrentSavings {
		months := int(math.Ceil((itemPrice - user.CurrentSavings) / surplus))
		if months < 1 {
			months = 1
		}
		wait = &months
	}

	if riskScore >= 6 {
		risk = "High"
		if wait != nil {
			recommendation := "High risk. This item costs more than your available savings. If you buy it now, your emergency fund will be too low. Wait about " + intString(*wait) + " months and save first."
			if impactGoals {
				recommendation += " This purchase may also affect your savings goals progress."
			}
			return models.AffordabilityResult{RiskLevel: risk, Reasons: reasons, Recommendation: recommendation, EstimatedWaitMonths: wait, Explanation: "Calculated using your income, expenses, savings, monthly surplus, emergency coverage, and goals."}
		}
		recommendation := "High risk. Do not buy this now. Build emergency savings first and reduce expenses before making this purchase."
		if surplus <= 0 {
			recommendation = "High risk. Your monthly surplus is zero or negative. Reduce expenses before considering this purchase."
		}
		return models.AffordabilityResult{RiskLevel: risk, Reasons: reasons, Recommendation: recommendation, EstimatedWaitMonths: wait, Explanation: "Calculated using your income, expenses, savings, monthly surplus, emergency coverage, and goals."}
	} else if riskScore >= 2 {
		risk = "Medium"
		recommendation := "Medium risk. This purchase needs caution. Save more first or reduce the price so your emergency savings stay protected."
		if wait != nil {
			recommendation = "Medium risk. Consider waiting about " + intString(*wait) + " months to save enough without affecting your emergency fund."
		}
		if impactGoals {
			recommendation += " This purchase may also affect your savings goals progress."
		}
		return models.AffordabilityResult{RiskLevel: risk, Reasons: reasons, Recommendation: recommendation, EstimatedWaitMonths: wait, Explanation: "Calculated using your income, expenses, savings, monthly surplus, emergency coverage, and goals."}
	}

	return models.AffordabilityResult{RiskLevel: risk, Reasons: reasons, Recommendation: "Low risk. You can likely afford this item without borrowing, and your emergency savings stay safe.", EstimatedWaitMonths: wait, Explanation: "Calculated using your income, expenses, savings, monthly surplus, emergency coverage, and goals."}
}

func round(v float64) float64 { return math.Round(v*100) / 100 }
func money(v float64) string  { return intString(int(math.Round(v))) }
func intString(i int) string {
	if i < 0 {
		return "-" + intString(-i)
	}
	s := ""
	for {
		s = string(rune('0'+i%10)) + s
		i /= 10
		if i == 0 {
			break
		}
	}
	return s
}
