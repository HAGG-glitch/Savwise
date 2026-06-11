package handlers

import "net/http"

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	uid, err := h.userIDFromRequest(r)
	if err != nil {
		respondError(w, 400, err.Error(), "missing_user_id")
		return
	}
	d, err := h.buildDashboard(uid)
	if err != nil {
		respondError(w, 500, "could not calculate dashboard", "dashboard_error")
		return
	}
	respondJSON(w, 200, "dashboard calculated", d)
}
