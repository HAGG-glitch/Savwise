package handlers

import (
	"net/http"
	"strings"

	"savwise-ai/internal/models"
)

func (h *Handler) Coach(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	var req models.CoachRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, 400, "invalid JSON", "invalid_json")
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		respondError(w, 400, "message is required", "validation_error")
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
		respondError(w, 403, "Please accept the prototype privacy notice before using Wizz.", "consent_required")
		return
	}
	dash, err := h.buildDashboard(req.UserID)
	if err != nil {
		respondError(w, 500, "could not prepare coaching context", "dashboard_error")
		return
	}
	out, err := h.Groq.Ask(req.Message, dash)
	if err != nil {
		respondError(w, 500, "could not generate coach response", "coach_error")
		return
	}
	displaySource := "Wizz"
	if out.Source == "groq" {
		displaySource = "Wizz (Groq)"
	}
	out.Source = displaySource
	_, _ = h.DB.Exec(`INSERT INTO ai_coach_messages(id,user_id,user_message,ai_response,model_name) VALUES($1,$2,$3,$4,$5)`, newUUID(), req.UserID, req.Message, out.Response, out.Model)
	respondJSON(w, 200, "coach response generated", out)
}
