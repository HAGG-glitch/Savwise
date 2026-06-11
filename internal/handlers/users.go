package handlers

import (
	"net/http"
	"strings"

	"savwise-ai/internal/models"
)

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req models.CreateUserRequest
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, 400, "invalid JSON", "invalid_json")
			return
		}
		if strings.TrimSpace(req.FullName) == "" {
			respondError(w, 400, "full name is required", "validation_error")
			return
		}
		if strings.TrimSpace(req.Email) == "" {
			respondError(w, 400, "email or username is required", "validation_error")
			return
		}
		if !req.ConsentAccepted {
			respondError(w, 400, "consent must be accepted", "consent_required")
			return
		}
		id := newUUID()
		_, err := h.DB.Exec(`INSERT INTO users(id, full_name, email, preferred_language) VALUES($1,$2,$3,$4)`,
			id, strings.TrimSpace(req.FullName), strings.TrimSpace(req.Email), "English")
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
				respondError(w, 409, "a user with this email already exists", "duplicate_email")
				return
			}
			respondError(w, 500, "could not create user", "user_create_failed")
			return
		}
		if req.MonthlyIncome > 0 || req.CurrentSavings > 0 || req.EmergencyTarget > 0 {
			_, err = h.DB.Exec(`INSERT INTO financial_profiles(id,user_id,monthly_income,current_savings,emergency_target,currency) VALUES($1,$2,$3,$4,$5,'SLE')`,
				newUUID(), id, req.MonthlyIncome, req.CurrentSavings, req.EmergencyTarget)
			if err != nil {
				respondError(w, 500, "could not create financial profile", "profile_create_failed")
				return
			}
		}
		_, _ = h.DB.Exec(`INSERT INTO consent_records(id,user_id,consent_type,accepted,policy_version) VALUES($1,$2,'prototype_data_processing',true,'1.0')`, newUUID(), id)
		if req.PreferredLanguage != "" {
			_, _ = h.DB.Exec(`UPDATE users SET preferred_language=$1 WHERE id=$2`, req.PreferredLanguage, id)
		}
		user, _ := h.getUser(id)
		respondJSON(w, 201, "user created", user)
	case http.MethodGet:
		uid, err := h.userIDFromRequest(r)
		if err != nil {
			respondError(w, 400, err.Error(), "missing_user_id")
			return
		}
		user, err := h.getUser(uid)
		if err != nil {
			respondError(w, 500, "could not load user", "user_error")
			return
		}
		respondJSON(w, 200, "user loaded", user)
	default:
		respondError(w, 405, "method not allowed", "method_not_allowed")
	}
}

func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	email := strings.TrimSpace(r.URL.Query().Get("email"))
	if email == "" {
		respondError(w, 400, "email query parameter is required", "missing_email")
		return
	}
	var id string
	err := h.DB.QueryRow(`SELECT id::text FROM users WHERE email=$1`, email).Scan(&id)
	if err != nil {
		respondError(w, 404, "user not found", "user_not_found")
		return
	}
	user, err := h.getUser(id)
	if err != nil {
		respondError(w, 500, "could not load user", "user_error")
		return
	}
	respondJSON(w, 200, "user found", user)
}

func (h *Handler) ResetAllDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	_, _ = h.DB.Exec(`DELETE FROM transactions; DELETE FROM savings_goals; DELETE FROM affordability_checks; DELETE FROM ai_coach_messages; DELETE FROM data_exports; DELETE FROM consent_records; DELETE FROM financial_profiles; DELETE FROM users;`)
	respondJSON(w, 200, "all demo users and data deleted", nil)
}
