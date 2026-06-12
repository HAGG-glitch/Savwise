package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"savwise-ai/internal/models"
)

type GroqService struct {
	APIKey string
	Model  string
	Client *http.Client
}

func NewGroqService(apiKey, model string) *GroqService {
	if model == "" {
		model = "llama-3.1-8b-instant"
	}
	return &GroqService{APIKey: apiKey, Model: model, Client: &http.Client{Timeout: 30 * time.Second}}
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}
type groqResponse struct {
	Choices []struct {
		Message groqMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (g *GroqService) Ask(message string, dashboard models.Dashboard) (models.CoachResponse, error) {
	if strings.TrimSpace(g.APIKey) == "" {
		return models.CoachResponse{Response: fallbackCoach(message), Source: "fallback", Model: "local-rules"}, nil
	}

	goalsSummary := ""
	if len(dashboard.Goals) > 0 {
		goalsSummary = " Active savings goals:"
		for _, goal := range dashboard.Goals {
			goalsSummary += fmt.Sprintf(" %s (target SLE %.0f, current SLE %.0f, progress %.0f%%),", goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.ProgressPercent)
		}
		goalsSummary = strings.TrimRight(goalsSummary, ",") + "."
	}

	spendingSummary := ""
	if len(dashboard.SpendingBreakdown) > 0 {
		spendingSummary = " Recent spending categories:"
		for _, cat := range dashboard.SpendingBreakdown {
			spendingSummary += fmt.Sprintf(" %s (SLE %.0f, %.0f%%),", cat.Category, cat.Amount, cat.Percent)
		}
		spendingSummary = strings.TrimRight(spendingSummary, ",") + "."
	}

	latestAffordability := ""
	if len(dashboard.Alerts) > 0 {
		latestAffordability = " Latest alerts available."
	}

	system := `You are Wizz, an educational financial coach for a university prototype in Sierra Leone. ONLY answer questions about personal finance, budgeting, saving, spending, osusu, financial literacy, and money management. If a question is NOT about personal finance (e.g. general knowledge, entertainment, politics, health advice, technology unrelated to finance, etc.), politely decline by saying "I am Wizz, your financial coach. I can only help with personal finance questions. Please ask me about budgeting, saving, or managing your money." Do not claim to be a bank, lender, lawyer, investment adviser, or licensed financial professional. Do not ask for mobile-money PINs, OTPs, bank passwords, national ID numbers, or confidential KYC information. Keep advice simple, safe, and explainable. Mention that recommendations are educational and should be verified for major financial decisions.

Use simple language. Keep responses short unless the user asks for detail. When useful, use this structure:
1. Quick answer
2. Why
3. What to do next
Use bullet points for steps. Use bold labels like **Risk**, **Reason**, **Next step**. Do not create one huge paragraph. If the user asks about buying an item, use their income, expenses, savings, goals, and recent spending context when available. If the user asks in Krio or about Krio terms, respond with simple Krio plus English explanation. Avoid complex financial jargon.`
	context := fmt.Sprintf("Financial context: monthly income SLE %.0f, monthly expenses SLE %.0f, monthly surplus SLE %.0f, savings rate %.1f%%, current savings SLE %.0f, emergency target SLE %.0f, health score %d/100.%s%s%s Do not reveal private identifiers.", dashboard.TotalIncome, dashboard.TotalExpenses, dashboard.MonthlySurplus, dashboard.SavingsRate, dashboard.User.CurrentSavings, dashboard.User.EmergencyTarget, dashboard.Score.Total, goalsSummary, spendingSummary, latestAffordability)
	payload := groqRequest{Model: g.Model, Temperature: 0.3, MaxTokens: 600, Messages: []groqMessage{{"system", system}, {"system", context}, {"user", message}}}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return models.CoachResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+g.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.Client.Do(req)
	if err != nil {
		return models.CoachResponse{Response: "Wizz is currently offline. Here is a rule-based educational suggestion instead: " + fallbackCoach(message), Source: "fallback", Model: "local-rules"}, nil
	}
	defer resp.Body.Close()
	var out groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return models.CoachResponse{}, err
	}
	if resp.StatusCode >= 300 || out.Error != nil {
		return models.CoachResponse{Response: "Wizz is currently offline. Here is a rule-based educational suggestion instead: " + fallbackCoach(message), Source: "fallback", Model: "local-rules"}, nil
	}
	if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
		return models.CoachResponse{}, errors.New("empty AI response")
	}
	return models.CoachResponse{Response: out.Choices[0].Message.Content, Source: "groq", Model: g.Model}, nil
}

var financeKeywords = []string{"save","spend","budget","money","income","expense","cost","price","afford","buy","goal","debt","loan","interest","emergency","fund","invest","osusu","salary","wage","earn","pay","bill","rent","food","transport","airtime","data","financial","savings","purchase","sell","trade","profit","bank","account","credit","cash","SLE","Leone"}

func isFinanceQuestion(m string) bool {
	mLower := strings.ToLower(m)
	for _, kw := range financeKeywords {
		if strings.Contains(mLower, kw) {
			return true
		}
	}
	return false
}

func fallbackCoach(message string) string {
	m := strings.ToLower(message)
	if !isFinanceQuestion(m) {
		return "I am Wizz, your financial coach. I can only help with personal finance questions. Please ask me about budgeting, saving, or managing your money. — Wizz"
	}
	switch {
	case strings.Contains(m, "osusu"):
		return "Osusu na group savings wey people dae contribute money regular, then one person dae collect the pot each round. Use am carefully: write who don pay, who don collect, and keep emergency money separate. This na educational guidance, not professional financial advice. — Wizz"
	case strings.Contains(m, "50") || strings.Contains(m, "budget"):
		return "Try the 50-30-20 rule: about 50% for needs, 30% for wants, and 20% for savings. Adjust am to your real income and essential costs for Salone. — Wizz"
	case strings.Contains(m, "trader") || strings.Contains(m, "business"):
		return "As a trader, separate business money from personal money. Record sales first, pay yourself a fixed amount, and leave restocking money untouched. — Wizz"
	case strings.Contains(m, "laptop") || strings.Contains(m, "buy") || strings.Contains(m, "afford"):
		return "Before you buy, check the price against your savings, income, monthly expenses, and emergency fund. If buying am go leave you with less than one month of expenses, wait small and save more first. — Wizz"
	default:
		return "Start small: track your income and expenses, set one emergency goal, and save a consistent amount each week. Even small savings reduce financial stress over time. — Wizz"
	}
}
