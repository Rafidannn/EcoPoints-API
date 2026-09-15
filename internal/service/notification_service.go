package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
)

type NotificationService interface {
	SendToUser(userID uint64, title, body string) error
}

type notificationService struct {
	projectID      string
	serviceAccount string
	tokenRepo      interface {
		GetTokensByUserID(userID uint64) ([]string, error)
	}
}

func NewNotificationService(
	serviceAccount string,
	tokenRepo interface {
		GetTokensByUserID(userID uint64) ([]string, error)
	},
) NotificationService {
	projectID := extractProjectID(serviceAccount)
	return &notificationService{
		projectID:      projectID,
		serviceAccount: serviceAccount,
		tokenRepo:      tokenRepo,
	}
}

func extractProjectID(serviceAccount string) string {
	if serviceAccount == "" {
		return ""
	}
	raw := serviceAccount
	if !strings.HasPrefix(strings.TrimSpace(serviceAccount), "{") {
		data, err := os.ReadFile(serviceAccount)
		if err != nil {
			return ""
		}
		raw = string(data)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return ""
	}
	if v, ok := m["project_id"].(string); ok {
		return v
	}
	return ""
}

func (s *notificationService) getAccessToken() (string, error) {
	raw := s.serviceAccount
	if !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		data, err := os.ReadFile(raw)
		if err != nil {
			return "", fmt.Errorf("baca service account: %w", err)
		}
		raw = string(data)
	}
	creds, err := google.CredentialsFromJSON(
		context.Background(),
		[]byte(raw),
		"https://www.googleapis.com/auth/firebase.messaging",
	)
	if err != nil {
		return "", fmt.Errorf("parse credentials: %w", err)
	}
	tok, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("ambil token: %w", err)
	}
	return tok.AccessToken, nil
}

func (s *notificationService) SendToUser(userID uint64, title, body string) error {
	if s.serviceAccount == "" || s.projectID == "" {
		return nil
	}

	tokens, err := s.tokenRepo.GetTokensByUserID(userID)
	if err != nil || len(tokens) == 0 {
		return nil
	}

	accessToken, err := s.getAccessToken()
	if err != nil {
		log.Printf("FCM access token error: %v", err)
		return nil
	}

	for _, token := range tokens {
		if err := s.sendFCM(accessToken, token, title, body); err != nil {
			log.Printf("FCM kirim gagal ke token %s: %v", token[:10], err)
		}
	}
	return nil
}

func (s *notificationService) sendFCM(accessToken, deviceToken, title, body string) error {
	url := fmt.Sprintf(
		"https://fcm.googleapis.com/v1/projects/%s/messages:send",
		s.projectID,
	)

	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"token": deviceToken,
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
			"android": map[string]interface{}{
				"notification": map[string]string{
					"channel_id": "ecopoints_channel",
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("FCM HTTP %d", resp.StatusCode)
	}
	return nil
}
