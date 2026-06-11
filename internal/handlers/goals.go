package handlers

import (
	"net/http"

	"savwise-ai/internal/models"
	"savwise-ai/internal/services"
)

func (h *Handler) Goals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uid, err := h.userIDFromRequest(r)
		if err != nil {
			respondError(w, 400, err.Error(), "missing_user_id")
			return
		}
		goals, err := h.listGoals(uid)
		if err != nil {
			respondError(w, 500, "could not list goals", "goal_list_failed")
			return
		}
		respondJSON(w, 200, "goals loaded", goals)
	case http.MethodPost:
		var req models.GoalRequest
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, 400, "invalid JSON", "invalid_json")
			return
		}
		if err := services.ValidateGoal(req); err != nil {
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
			respondError(w, 403, "Please accept the prototype privacy notice before creating goals.", "consent_required")
			return
		}
		_, err = h.DB.Exec(`INSERT INTO savings_goals(id,user_id,name,target_amount,current_amount,monthly_contribution) VALUES($1,$2,$3,$4,$5,$6)`, newUUID(), req.UserID, req.Name, req.TargetAmount, req.CurrentAmount, req.MonthlyContribution)
		if err != nil {
			respondError(w, 500, "could not save goal", "goal_save_failed")
			return
		}
		goals, _ := h.listGoals(req.UserID)
		respondJSON(w, 201, "goal saved", goals)
	default:
		respondError(w, 405, "method not allowed", "method_not_allowed")
	}
}

func (h *Handler) GoalByID(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userIDFromRequest(r)
	if err != nil {
		respondError(w, 400, err.Error(), "missing_user_id")
		return
	}
	id := idFromPath(r.URL.Path, "/api/goals/")
	if id == "" {
		respondError(w, 400, "goal id is required", "missing_id")
		return
	}
	switch r.Method {
	case http.MethodPut:
		var req models.GoalRequest
		if err := decodeJSON(r, &req); err != nil {
			respondError(w, 400, "invalid JSON", "invalid_json")
			return
		}
		if err := services.ValidateGoal(req); err != nil {
			respondError(w, 400, err.Error(), "validation_error")
			return
		}
		res, err := h.DB.Exec(`UPDATE savings_goals SET name=$1,target_amount=$2,current_amount=$3,monthly_contribution=$4,updated_at=NOW() WHERE id=$5 AND user_id=$6`, req.Name, req.TargetAmount, req.CurrentAmount, req.MonthlyContribution, id, uid)
		if err != nil {
			respondError(w, 500, "could not update goal", "goal_update_failed")
			return
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			respondError(w, 404, "goal not found", "not_found")
			return
		}
		goals, _ := h.listGoals(uid)
		respondJSON(w, 200, "goal updated", goals)
	case http.MethodDelete:
		res, err := h.DB.Exec(`DELETE FROM savings_goals WHERE id=$1 AND user_id=$2`, id, uid)
		if err != nil {
			respondError(w, 500, "could not delete goal", "goal_delete_failed")
			return
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			respondError(w, 404, "goal not found", "not_found")
			return
		}
		respondJSON(w, 200, "goal deleted", nil)
	default:
		respondError(w, 405, "method not allowed", "method_not_allowed")
	}
}
