package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"savwise-ai/internal/config"
	"savwise-ai/internal/models"
	"savwise-ai/internal/services"
)

type Handler struct {
	DB   *sql.DB
	Cfg  config.Config
	Groq *services.GroqService
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/users", h.Users)
	mux.HandleFunc("/api/current-user", h.CurrentUser)
	mux.HandleFunc("/api/profile", h.Profile)
	mux.HandleFunc("/api/consent", h.Consent)
	mux.HandleFunc("/api/transactions", h.Transactions)
	mux.HandleFunc("/api/transactions/", h.TransactionByID)
	mux.HandleFunc("/api/goals", h.Goals)
	mux.HandleFunc("/api/goals/", h.GoalByID)
	mux.HandleFunc("/api/dashboard", h.Dashboard)
	mux.HandleFunc("/api/affordability", h.Affordability)
	mux.HandleFunc("/api/coach", h.Coach)
	mux.HandleFunc("/api/export/json", h.ExportJSON)
	mux.HandleFunc("/api/export/csv", h.ExportCSV)
	mux.HandleFunc("/api/import/json", h.ImportJSON)
	mux.HandleFunc("/api/import/csv", h.ImportCSV)
	mux.HandleFunc("/api/load-demo", h.LoadDemo)
	mux.HandleFunc("/api/reset", h.Reset)
	mux.HandleFunc("/api/reset-all-demo-data", h.ResetAllDemo)
}

func respondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.APIResponse{Success: status < 400, Message: message, Data: data})
}

func respondError(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.APIResponse{Success: false, Message: message, Error: code})
}

func decodeJSON(r *http.Request, v interface{}) error { return json.NewDecoder(r.Body).Decode(v) }

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:32]
}

func idFromPath(path, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(path, prefix))
}

func (h *Handler) userIDFromRequest(r *http.Request) (string, error) {
	uid := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if uid == "" {
		return "", errors.New("user_id query parameter is required")
	}
	var exists bool
	err := h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, uid).Scan(&exists)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("user not found")
	}
	return uid, nil
}

func (h *Handler) getUser(uid string) (models.User, error) {
	var u models.User
	var consent bool
	err := h.DB.QueryRow(`
        SELECT u.id::text, u.full_name, u.email, u.preferred_language,
               COALESCE(fp.monthly_income,0), COALESCE(fp.current_savings,0), COALESCE(fp.emergency_target,1500),
               EXISTS(SELECT 1 FROM consent_records c WHERE c.user_id=u.id AND c.accepted=true)
        FROM users u
        LEFT JOIN financial_profiles fp ON fp.user_id=u.id
        WHERE u.id=$1`, uid).Scan(&u.ID, &u.FullName, &u.Email, &u.PreferredLanguage, &u.MonthlyIncome, &u.CurrentSavings, &u.EmergencyTarget, &consent)
	u.ConsentAccepted = consent
	return u, err
}

func (h *Handler) listTransactions(uid string) ([]models.Transaction, error) {
	rows, err := h.DB.Query(`SELECT id::text, user_id::text, description, amount::float8, type, category, transaction_date::text, created_at::text FROM transactions WHERE user_id=$1 ORDER BY transaction_date DESC, created_at DESC`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.Description, &tx.Amount, &tx.Type, &tx.Category, &tx.Date, &tx.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, tx)
	}
	return out, rows.Err()
}

func (h *Handler) listGoals(uid string) ([]models.Goal, error) {
	rows, err := h.DB.Query(`SELECT id::text, user_id::text, name, target_amount::float8, current_amount::float8, monthly_contribution::float8, created_at::text FROM savings_goals WHERE user_id=$1 ORDER BY created_at DESC`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Goal
	for rows.Next() {
		var g models.Goal
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.MonthlyContribution, &g.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return services.EnhanceGoals(out), rows.Err()
}

func (h *Handler) buildDashboard(uid string) (models.Dashboard, error) {
	u, err := h.getUser(uid)
	if err != nil {
		return models.Dashboard{}, err
	}
	txs, err := h.listTransactions(uid)
	if err != nil {
		return models.Dashboard{}, err
	}
	goals, err := h.listGoals(uid)
	if err != nil {
		return models.Dashboard{}, err
	}
	return services.CalculateDashboard(u, txs, goals), nil
}
