package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"savwise-ai/internal/models"
)

func TransactionsCSV(transactions []models.Transaction) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"id", "description", "amount", "type", "category", "date", "createdAt"}); err != nil {
		return nil, err
	}
	for _, tx := range transactions {
		if err := writer.Write([]string{tx.ID, tx.Description, fmt.Sprintf("%.2f", tx.Amount), tx.Type, tx.Category, tx.Date, tx.CreatedAt}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

type CSVImportResult struct {
	TotalRows     int      `json:"totalRows"`
	Transactions  int      `json:"transactions"`
	TotalIncome   float64  `json:"totalIncome"`
	TotalExpenses float64  `json:"totalExpenses"`
	Errors        []string `json:"errors"`
}

func ParseTransactionsCSV(r io.Reader) ([]models.Transaction, CSVImportResult, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, CSVImportResult{}, fmt.Errorf("could not read CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, CSVImportResult{}, fmt.Errorf("CSV must have a header row and at least one data row")
	}

	header := records[0]
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	required := []string{"description", "amount", "type", "category", "date"}
	for _, r := range required {
		if _, ok := headerMap[r]; !ok {
			return nil, CSVImportResult{}, fmt.Errorf("missing required CSV column: %s", r)
		}
	}

	var txs []models.Transaction
	var errors []string
	var totalIncome, totalExpenses float64

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < len(required) {
			errors = append(errors, fmt.Sprintf("row %d: too few columns", i+1))
			continue
		}

		desc := strings.TrimSpace(row[headerMap["description"]])
		if desc == "" {
			errors = append(errors, fmt.Sprintf("row %d: description is required", i+1))
			continue
		}

		amountStr := strings.TrimSpace(row[headerMap["amount"]])
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			errors = append(errors, fmt.Sprintf("row %d: invalid amount '%s'", i+1, amountStr))
			continue
		}

		txType := strings.TrimSpace(strings.ToLower(row[headerMap["type"]]))
		if txType != "income" && txType != "expense" {
			errors = append(errors, fmt.Sprintf("row %d: type must be 'income' or 'expense', got '%s'", i+1, txType))
			continue
		}

		category := strings.TrimSpace(row[headerMap["category"]])
		if category == "" {
			errors = append(errors, fmt.Sprintf("row %d: category is required", i+1))
			continue
		}

		date := strings.TrimSpace(row[headerMap["date"]])
		if _, err := time.Parse("2006-01-02", date); err != nil {
			errors = append(errors, fmt.Sprintf("row %d: invalid date '%s', use YYYY-MM-DD", i+1, date))
			continue
		}

		tx := models.Transaction{
			Description: desc,
			Amount:      amount,
			Type:        txType,
			Category:    category,
			Date:        date,
		}

		if idIdx, ok := headerMap["id"]; ok && idIdx < len(row) {
			tx.ID = strings.TrimSpace(row[idIdx])
		}

		if txType == "income" {
			totalIncome += amount
		} else {
			totalExpenses += amount
		}
		txs = append(txs, tx)
	}

	result := CSVImportResult{
		TotalRows:     len(records) - 1,
		Transactions:  len(txs),
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		Errors:        errors,
	}
	return txs, result, nil
}
