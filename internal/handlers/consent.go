package handlers

import "net/http"

func (h *Handler) Consent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	uid := r.URL.Query().Get("user_id")
	if uid == "" {
		respondError(w, 400, "user_id is required", "missing_user_id")
		return
	}
	_, err := h.DB.Exec(`INSERT INTO consent_records(id,user_id,consent_type,accepted,policy_version) VALUES($1,$2,'prototype_data_processing',true,'1.0')`, newUUID(), uid)
	if err != nil {
		respondError(w, 500, "could not save consent", "consent_error")
		return
	}
	respondJSON(w, 200, "consent saved", map[string]bool{"consentAccepted": true})
}
