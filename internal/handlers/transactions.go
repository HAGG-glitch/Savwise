package handlers

import (
	"net/http"

	"savwise-ai/internal/models"
	"savwise-ai/internal/services"
)

func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uid, err := h.userIDFromRequest(r)
		if err != nil {
			respondError(w, 400, err.Error(), "missing_user_id")
			return
		}
		txs, err := h.listTransactions(uid)
		if err != nil {
			respondError(w, 500, "could not list transactions", "transaction_list_failed")
			return
		}
		respondJSON(w, 200, "transactions loaded", txs)
	case http.MethodPost:
		var req models.TransactionRequest
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, 400, "invalid JSON", "invalid_json")
			return
		}
		if err := services.ValidateTransaction(req); err != nil {
			respondError(w, 400, err.Error(), "validation_error")
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
			respondError(w, 403, "Please accept the prototype privacy notice before adding transactions.", "consent_required")
			return
		}
		id := newUUID()
		_, err = h.DB.Exec(`INSERT INTO transactions(id,user_id,description,amount,type,category,transaction_date) VALUES($1,$2,$3,$4,$5,$6,$7)`, id, req.UserID, req.Description, req.Amount, req.Type, req.Category, req.Date)
		if err != nil {
			respondError(w, 500, "could not save transaction", "transaction_save_failed")
			return
		}
		txs, _ := h.listTransactions(req.UserID)
		respondJSON(w, 201, "transaction saved", txs)
	default:
		respondError(w, 405, "method not allowed", "method_not_allowed")
	}
}

func (h *Handler) TransactionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, 405, "method not allowed", "method_not_allowed")
		return
	}
	uid, err := h.userIDFromRequest(r)
	if err != nil {
		respondError(w, 400, err.Error(), "missing_user_id")
		return
	}
	id := idFromPath(r.URL.Path, "/api/transactions/")
	if id == "" {
		respondError(w, 400, "transaction id is required", "missing_id")
		return
	}
	res, err := h.DB.Exec(`DELETE FROM transactions WHERE id=$1 AND user_id=$2`, id, uid)
	if err != nil {
		respondError(w, 500, "could not delete transaction", "transaction_delete_failed")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		respondError(w, 404, "transaction not found", "not_found")
		return
	}
	respondJSON(w, 200, "transaction deleted", nil)
}
