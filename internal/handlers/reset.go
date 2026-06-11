package handlers

import "net/http"

func (h *Handler) Reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	uid, err := h.userIDFromRequest(r)
	if err != nil {
		respondError(w, 400, err.Error(), "missing_user_id")
		return
	}
	_, err = h.DB.Exec(`DELETE FROM transactions WHERE user_id=$1; DELETE FROM savings_goals WHERE user_id=$1; DELETE FROM affordability_checks WHERE user_id=$1; DELETE FROM ai_coach_messages WHERE user_id=$1; DELETE FROM data_exports WHERE user_id=$1;`, uid)
	if err != nil {
		respondError(w, 500, "could not reset demo data", "reset_error")
		return
	}
	respondJSON(w, 200, "demo data reset for user", nil)
}

func (h *Handler) LoadDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	uid := r.URL.Query().Get("user_id")
	if uid == "" {
		respondError(w, 400, "user_id is required", "missing_user_id")
		return
	}
	tx, err := h.DB.Begin()
	if err != nil {
		respondError(w, 500, "could not start demo load", "database_error")
		return
	}
	defer tx.Rollback()
	_, _ = tx.Exec(`DELETE FROM transactions WHERE user_id=$1`, uid)
	_, _ = tx.Exec(`DELETE FROM savings_goals WHERE user_id=$1`, uid)
	_, _ = tx.Exec(`UPDATE users SET full_name='Fatmata Kamara', preferred_language='English' WHERE id=$1`, uid)
	_, err = tx.Exec(`INSERT INTO financial_profiles(id,user_id,monthly_income,current_savings,emergency_target,currency) VALUES($1,$2,1200,450,1500,'SLE') ON CONFLICT(user_id) DO UPDATE SET monthly_income=1200,current_savings=450,emergency_target=1500,updated_at=NOW()`, newUUID(), uid)
	if err != nil {
		respondError(w, 500, "could not load profile", "demo_error")
		return
	}
	_, _ = tx.Exec(`INSERT INTO consent_records(id,user_id,consent_type,accepted,policy_version) VALUES($1,$2,'prototype_data_processing',true,'1.0')`, newUUID(), uid)
	demoTx := []struct {
		desc           string
		amount         float64
		typ, cat, date string
	}{
		{"Monthly allowance", 1200, "income", "Income", "2026-06-01"},
		{"Food", 400, "expense", "Food", "2026-06-02"},
		{"Education", 300, "expense", "Education", "2026-06-03"},
		{"Transport", 200, "expense", "Transport", "2026-06-04"},
		{"Airtime/Data", 180, "expense", "Airtime/Data", "2026-06-05"},
		{"Entertainment", 120, "expense", "Entertainment", "2026-06-06"},
	}
	for _, d := range demoTx {
		_, err = tx.Exec(`INSERT INTO transactions(id,user_id,description,amount,type,category,transaction_date) VALUES($1,$2,$3,$4,$5,$6,$7)`, newUUID(), uid, d.desc, d.amount, d.typ, d.cat, d.date)
		if err != nil {
			respondError(w, 500, "could not load transactions", "demo_error")
			return
		}
	}
	_, err = tx.Exec(`INSERT INTO savings_goals(id,user_id,name,target_amount,current_amount,monthly_contribution) VALUES($1,$2,'Emergency Fund',1500,450,200),($3,$2,'Textbook Fund',1000,680,160)`, newUUID(), uid, newUUID())
	if err != nil {
		respondError(w, 500, "could not load goals", "demo_error")
		return
	}
	if err := tx.Commit(); err != nil {
		respondError(w, 500, "could not finish demo load", "demo_error")
		return
	}
	respondJSON(w, 200, "demo data loaded for user", nil)
}
