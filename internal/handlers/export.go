package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"savwise-ai/internal/models"
	"savwise-ai/internal/services"
)

func (h *Handler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
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
	if !user.ConsentAccepted {
		respondError(w, 403, "Please accept the prototype privacy notice before exporting data.", "consent_required")
		return
	}
	txs, _ := h.listTransactions(uid)
	goals, _ := h.listGoals(uid)
	pkg := models.ExportPackage{SchemaVersion: "1.0", ExportedAt: time.Now().Format(time.RFC3339), User: user, Transactions: txs, Goals: goals}
	_, _ = h.DB.Exec(`INSERT INTO data_exports(id,user_id,export_format) VALUES($1,$2,'json')`, newUUID(), uid)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=savwise-export.json")
	_ = json.NewEncoder(w).Encode(pkg)
}

func (h *Handler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
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
	if !user.ConsentAccepted {
		respondError(w, 403, "Please accept the prototype privacy notice before exporting data.", "consent_required")
		return
	}
	txs, _ := h.listTransactions(uid)
	data, err := services.TransactionsCSV(txs)
	if err != nil {
		respondError(w, 500, "could not create CSV", "csv_error")
		return
	}
	_, _ = h.DB.Exec(`INSERT INTO data_exports(id,user_id,export_format) VALUES($1,$2,'csv')`, newUUID(), uid)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=savwise-transactions.csv")
	_, _ = w.Write(data)
}

func (h *Handler) ImportJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	var pkg models.ExportPackage
	if err := decodeJSON(r, &pkg); err != nil {
		respondError(w, 400, "invalid JSON import", "invalid_json")
		return
	}
	if pkg.SchemaVersion != "1.0" {
		respondError(w, 400, "unsupported schema version", "schema_error")
		return
	}
	uid := r.URL.Query().Get("user_id")
	if uid == "" {
		respondError(w, 400, "user_id is required", "missing_user_id")
		return
	}
	user, err := h.getUser(uid)
	if err != nil {
		respondError(w, 500, "could not load user", "user_error")
		return
	}
	if !user.ConsentAccepted {
		respondError(w, 403, "Please accept the prototype privacy notice before importing data.", "consent_required")
		return
	}
	tx, err := h.DB.Begin()
	if err != nil {
		respondError(w, 500, "could not start import", "database_error")
		return
	}
	defer tx.Rollback()
	_, _ = tx.Exec(`DELETE FROM transactions WHERE user_id=$1`, uid)
	_, _ = tx.Exec(`DELETE FROM savings_goals WHERE user_id=$1`, uid)
	_, err = tx.Exec(`UPDATE users SET full_name=$1, email=$2, preferred_language=$3 WHERE id=$4`, pkg.User.FullName, pkg.User.Email, pkg.User.PreferredLanguage, uid)
	if err != nil {
		respondError(w, 500, "could not import user", "import_error")
		return
	}
	_, err = tx.Exec(`INSERT INTO financial_profiles(id,user_id,monthly_income,current_savings,emergency_target,currency) VALUES($1,$2,$3,$4,$5,'SLE') ON CONFLICT(user_id) DO UPDATE SET monthly_income=EXCLUDED.monthly_income,current_savings=EXCLUDED.current_savings,emergency_target=EXCLUDED.emergency_target,updated_at=NOW()`, newUUID(), uid, pkg.User.MonthlyIncome, pkg.User.CurrentSavings, pkg.User.EmergencyTarget)
	if err != nil {
		respondError(w, 500, "could not import profile", "import_error")
		return
	}
	for _, t := range pkg.Transactions {
		if t.ID == "" {
			t.ID = newUUID()
		}
		_, err = tx.Exec(`INSERT INTO transactions(id,user_id,description,amount,type,category,transaction_date) VALUES($1,$2,$3,$4,$5,$6,$7)`, t.ID, uid, t.Description, t.Amount, t.Type, t.Category, t.Date)
		if err != nil {
			respondError(w, 400, "invalid transaction in import", "import_error")
			return
		}
	}
	for _, g := range pkg.Goals {
		if g.ID == "" {
			g.ID = newUUID()
		}
		_, err = tx.Exec(`INSERT INTO savings_goals(id,user_id,name,target_amount,current_amount,monthly_contribution) VALUES($1,$2,$3,$4,$5,$6)`, g.ID, uid, g.Name, g.TargetAmount, g.CurrentAmount, g.MonthlyContribution)
		if err != nil {
			respondError(w, 400, "invalid goal in import", "import_error")
			return
		}
	}
	if err := tx.Commit(); err != nil {
		respondError(w, 500, "could not complete import", "import_error")
		return
	}
	respondJSON(w, 200, "JSON backup imported", nil)
}

func (h *Handler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	uid := r.URL.Query().Get("user_id")
	if uid == "" {
		respondError(w, 400, "user_id is required", "missing_user_id")
		return
	}
	user, err := h.getUser(uid)
	if err != nil {
		respondError(w, 500, "could not load user", "user_error")
		return
	}
	if !user.ConsentAccepted {
		respondError(w, 403, "Please accept the prototype privacy notice before importing data.", "consent_required")
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, 400, "CSV file is required", "missing_file")
		return
	}
	defer file.Close()

	txs, result, err := services.ParseTransactionsCSV(file)
	if err != nil {
		respondError(w, 400, err.Error(), "csv_parse_error")
		return
	}

	if len(result.Errors) > 0 {
		respondJSON(w, 200, "CSV parsed with errors", result)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		respondError(w, 500, "could not start import", "database_error")
		return
	}
	defer tx.Rollback()

	for _, t := range txs {
		if t.ID == "" {
			t.ID = newUUID()
		}
		_, err = tx.Exec(`INSERT INTO transactions(id,user_id,description,amount,type,category,transaction_date) VALUES($1,$2,$3,$4,$5,$6,$7)`, t.ID, uid, t.Description, t.Amount, t.Type, t.Category, t.Date)
		if err != nil {
			respondError(w, 400, "could not import CSV row: "+t.Description, "import_error")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondError(w, 500, "could not complete CSV import", "import_error")
		return
	}

	respondJSON(w, 200, "CSV imported successfully", result)
}
