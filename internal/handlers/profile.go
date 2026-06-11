package handlers

import (
	"net/http"

	"savwise-ai/internal/models"
	"savwise-ai/internal/services"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	if err := h.DB.Ping(); err != nil {
		respondError(w, 500, "database is not reachable", "database_error")
		return
	}
	respondJSON(w, 200, "SavWise AI API is healthy", map[string]string{"status": "ok"})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uid, err := h.userIDFromRequest(r)
		if err != nil {
			respondError(w, 400, err.Error(), "missing_user_id")
			return
		}
		user, err := h.getUser(uid)
		if err != nil {
			respondError(w, 500, "could not load profile", "profile_load_failed")
			return
		}
		respondJSON(w, 200, "profile loaded", user)
	case http.MethodPost:
		var req models.ProfileRequest
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, 400, "invalid JSON", "invalid_json")
			return
		}
		if err := services.ValidateProfile(req); err != nil {
			respondError(w, 400, err.Error(), "validation_error")
			return
		}
		if req.UserID == "" {
			respondError(w, 400, "user_id is required", "missing_user_id")
			return
		}
		if req.PreferredLanguage == "" {
			req.PreferredLanguage = "English"
		}
		_, err := h.DB.Exec(`UPDATE users SET full_name=$1, email=$2, preferred_language=$3 WHERE id=$4`,
			req.FullName, req.Email, req.PreferredLanguage, req.UserID)
		if err != nil {
			respondError(w, 500, "could not update user", "profile_update_failed")
			return
		}
		_, err = h.DB.Exec(`
            INSERT INTO financial_profiles(id,user_id,monthly_income,current_savings,emergency_target,currency)
            VALUES($1,$2,$3,$4,$5,'SLE')
            ON CONFLICT(user_id) DO UPDATE SET monthly_income=EXCLUDED.monthly_income, current_savings=EXCLUDED.current_savings, emergency_target=EXCLUDED.emergency_target, updated_at=NOW()`,
			newUUID(), req.UserID, req.MonthlyIncome, req.CurrentSavings, req.EmergencyTarget)
		if err != nil {
			respondError(w, 500, "could not update financial profile", "profile_update_failed")
			return
		}
		_, _ = h.DB.Exec(`INSERT INTO consent_records(id,user_id,consent_type,accepted,policy_version) VALUES($1,$2,'prototype_data_processing',true,'1.0')`, newUUID(), req.UserID)
		user, _ := h.getUser(req.UserID)
		respondJSON(w, 200, "profile saved", user)
	default:
		respondError(w, 405, "method not allowed", "method_not_allowed")
	}
}
