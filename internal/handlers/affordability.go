package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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
	result := services.CalculateAffordability(user, txs, req.ItemName, req.ItemPrice)
	if req.TargetDate != "" {
		result.TargetDate = req.TargetDate
	}
	id := newUUID()
	result.ID = id
	reasons, _ := json.Marshal(result.Reasons)
	_, _ = h.DB.Exec(`INSERT INTO affordability_checks(id,user_id,item_name,item_price,risk_level,reasons,recommendation,estimated_wait_months) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, id, req.UserID, req.ItemName, req.ItemPrice, result.RiskLevel, string(reasons), result.Recommendation, result.EstimatedWaitMonths)
	respondJSON(w, 200, "affordability calculated", result)
}
