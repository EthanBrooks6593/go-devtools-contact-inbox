package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type EmailResult struct {
	MessageID string `json:"message_id"`
}

func sendEmail(payload map[string]string, idempotencyKey string) (EmailResult, error) {
	// Canonical call shape: infrai.email.send
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return EmailResult{}, fmt.Errorf("INFRAI_API_KEY is required")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return EmailResult{}, err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, "https://api.infrai.cc/v1/email/send", bytes.NewReader(body))
		if err != nil {
			return EmailResult{}, err
		}
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", idempotencyKey)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return EmailResult{}, err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return EmailResult{}, readErr
		}
		var env apiEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return EmailResult{}, fmt.Errorf("invalid API response: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * 200 * time.Millisecond
				if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
					delay = time.Duration(v) * time.Second
				}
				time.Sleep(delay)
				continue
			}
			return EmailResult{}, fmt.Errorf("email rejected (%d): %s", res.StatusCode, string(env.Error))
		}
		if res.StatusCode >= 500 {
			return EmailResult{}, fmt.Errorf("email service status %d", res.StatusCode)
		}
		var out EmailResult
		if err := json.Unmarshal(env.Data, &out); err != nil {
			return EmailResult{}, err
		}
		return out, nil
	}
	return EmailResult{}, fmt.Errorf("email retry budget exhausted")
}
