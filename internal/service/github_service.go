package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yashashkumar/ai-issue-agent/internal/errors"
)

type GitHubIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	User   struct {
		Login string `json:"login"`
		Email string `json:"email"` // Often empty unless specifically configured
	} `json:"user"`
}

type GitHubComment struct {
	Body string `json:"body"`
	User struct {
		Login string `json:"login"`
		Email string `json:"email"`
	} `json:"user"`
}

type GitHubWebhookPayload struct {
	Action  string         `json:"action"`
	Issue   GitHubIssue    `json:"issue"`
	Comment *GitHubComment `json:"comment"`
	Sender  struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func VerifySignature(secret string, signatureHeader string, body []byte) error {
	const signaturePrefix = "sha256="
	if !strings.HasPrefix(signatureHeader, signaturePrefix) {
		return errors.Wrap(errors.ErrUnauthorized, "invalid signature format")
	}

	actualMAC, err := hex.DecodeString(signatureHeader[len(signaturePrefix):])
	if err != nil {
		return errors.Wrap(errors.ErrUnauthorized, "failed to decode signature")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)

	if subtle.ConstantTimeCompare(actualMAC, expectedMAC) != 1 {
		return errors.Wrap(errors.ErrUnauthorized, "signature mismatch")
	}

	return nil
}

func ParseWebhookPayload(body []byte) (*GitHubWebhookPayload, error) {
	var payload GitHubWebhookPayload
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode json payload: %w", err)
	}
	return &payload, nil
}
