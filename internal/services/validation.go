package services

import (
	"errors"
	"strings"
	"time"

	"savwise-ai/internal/models"
)

func ValidateTransaction(t models.TransactionRequest) error {
	if strings.TrimSpace(t.Description) == "" {
		return errors.New("description is required")
	}
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if t.Type != "income" && t.Type != "expense" {
		return errors.New("type must be income or expense")
	}
	if strings.TrimSpace(t.Category) == "" {
		return errors.New("category is required")
	}
	if _, err := time.Parse("2006-01-02", t.Date); err != nil {
		return errors.New("date must use YYYY-MM-DD format")
	}
	return nil
}

func ValidateGoal(g models.GoalRequest) error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("goal name is required")
	}
	if g.TargetAmount <= 0 {
		return errors.New("target amount must be greater than zero")
	}
	if g.CurrentAmount < 0 || g.MonthlyContribution < 0 {
		return errors.New("current amount and monthly contribution cannot be negative")
	}
	return nil
}

func ValidateProfile(p models.ProfileRequest) error {
	if strings.TrimSpace(p.FullName) == "" {
		return errors.New("full name is required")
	}
	if p.MonthlyIncome < 0 || p.CurrentSavings < 0 || p.EmergencyTarget < 0 {
		return errors.New("financial values cannot be negative")
	}
	if !p.ConsentAccepted {
		return errors.New("prototype consent must be accepted before saving profile data")
	}
	return nil
}
