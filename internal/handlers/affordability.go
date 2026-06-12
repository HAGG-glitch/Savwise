package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"savwise-ai/internal/models"
	"savwise-ai/internal/services"
)

func (h *Handler) Affordability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	var req models.AffordabilityRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, 400, "invalid JSON", "invalid_json")
		return
	}
	if strings.TrimSpace(req.ItemName) == "" || req.ItemPrice <= 0 {
		respondError(w, 400, "item name and positive item price are required", "validation_error")
		return
	}
	if req.UserID == "" {
		respondError(w, 400, "user_id is required", "missing_user_id")
		return
	}
	if req.TargetDate != "" {
		if _, err := time.Parse("2006-01-02", req.TargetDate); err != nil {
			respondError(w, 400, "target date must use YYYY-MM-DD format and be a valid date", "validation_error")
			return
		}
		parsed, _ := time.Parse("2006-01-02", req.TargetDate)
		today := time.Now().Truncate(24 * time.Hour)
		if parsed.Before(today) {
			respondError(w, 400, "target date cannot be before today", "validation_error")
			return
		}
	}
	user, err := h.getUser(req.UserID)
	if err != nil {
		respondError(w, 500, "could not load user", "user_error")
		return
	}
	if !user.ConsentAccepted {
		respondError(w, 403, "Please accept the prototype privacy notice before using the affordability checker.", "consent_required")
		return
	}
	txs, err := h.listTransactions(req.UserID)
	if err != nil {
		respondError(w, 500, "could not load transactions", "transaction_error")
		return
	}
	goals, err := h.listGoals(req.UserID)
	if err != nil {
		respondError(w, 500, "could not load goals", "goal_error")
		return
	}
	result := services.CalculateAffordability(user, txs, goals, req.ItemName, req.ItemPrice, req.TargetDate)
	id := newUUID()
	result.ID = id
	reasons, _ := json.Marshal(result.Reasons)
	_, _ = h.DB.Exec(`INSERT INTO affordability_checks(id,user_id,item_name,item_price,risk_level,reasons,recommendation,estimated_wait_months,target_date,calculated_at,monthly_income,monthly_expenses,monthly_surplus,funding_gap,months_until_target,required_monthly_saving,active_goal_commitments,available_after_goals,expense_period) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`,
		id, req.UserID, req.ItemName, req.ItemPrice, result.RiskLevel, string(reasons), result.Recommendation, nil, result.TargetDate, result.CalculatedAt, result.MonthlyIncome, result.MonthlyExpenses, result.MonthlySurplus, result.FundingGap, result.MonthsUntilTarget, result.RequiredMonthlySaving, result.ActiveGoalCommitments, result.AvailableAfterGoals, result.ExpensePeriod)
	respondJSON(w, 200, "affordability calculated", result)
}
