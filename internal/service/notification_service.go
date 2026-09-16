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
	SendToAll(title, body string) (int, error)
}

type notificationService struct {
	projectID      string
	serviceAccount string
	tokenRepo      interface {
		GetTokensByUserID(userID uint64) ([]string, error)
		GetAllTokens() ([]string, error)
	}
}

func NewNotificationService(
	serviceAccount string,
	tokenRepo interface {
		GetTokensByUserID(userID uint64) ([]string, error)
		GetAllTokens() ([]string, error)
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
	raw := serviceAccount
	if raw == "" {
		raw = "serviceAccountKey.json"
	}
	if !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		data, err := os.ReadFile(raw)
		if err != nil {
			data, err = os.ReadFile("serviceAccountKey.json")
			if err != nil {
				log.Printf("FCM extractProjectID error membaca %s: %v", raw, err)
				return ""
			}
		}
		raw = string(data)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		log.Printf("FCM extractProjectID unmarshal error: %v", err)
		return ""
	}
	if v, ok := m["project_id"].(string); ok {
		return v
	}
	return ""
}

func (s *notificationService) getAccessToken() (string, error) {
	raw := s.serviceAccount
	if raw == "" {
		raw = "serviceAccountKey.json"
	}
	if !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		data, err := os.ReadFile(raw)
		if err != nil {
			data, err = os.ReadFile("serviceAccountKey.json")
			if err != nil {
				return "", fmt.Errorf("file Firebase Service Account tidak ditemukan di '%s' maupun 'serviceAccountKey.json': %w", s.serviceAccount, err)
			}
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
		return "", fmt.Errorf("ambil token OAuth: %w", err)
	}
	return tok.AccessToken, nil
}

func (s *notificationService) SendToUser(userID uint64, title, body string) error {
	if s.projectID == "" {
		s.projectID = extractProjectID(s.serviceAccount)
	}
	if s.projectID == "" {
		return fmt.Errorf("Firebase Project ID tidak ditemukan di serviceAccountKey.json")
	}

	tokens, err := s.tokenRepo.GetTokensByUserID(userID)
	if err != nil {
		return fmt.Errorf("gagal mengambil token perangkat user: %w", err)
	}
	if len(tokens) == 0 {
		return fmt.Errorf("nasabah ini belum memiliki perangkat smartphone yang terhubung (FCM Token belum terdaftar)")
	}

	accessToken, err := s.getAccessToken()
	if err != nil {
		log.Printf("FCM access token error: %v", err)
		return fmt.Errorf("gagal otentikasi Firebase FCM: %w", err)
	}

	var lastErr error
	sentCount := 0
	for _, token := range tokens {
		if err := s.sendFCM(accessToken, token, title, body); err != nil {
			log.Printf("FCM kirim gagal ke token %s: %v", token, err)
			lastErr = err
		} else {
			sentCount++
		}
	}

	if sentCount == 0 && lastErr != nil {
		return fmt.Errorf("gagal mengirim ke FCM Google: %w", lastErr)
	}
	return nil
}

func (s *notificationService) SendToAll(title, body string) (int, error) {
	if s.projectID == "" {
		s.projectID = extractProjectID(s.serviceAccount)
	}
	if s.projectID == "" {
		return 0, fmt.Errorf("Firebase Project ID tidak ditemukan di serviceAccountKey.json")
	}

	tokens, err := s.tokenRepo.GetAllTokens()
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil daftar token perangkat: %w", err)
	}
	if len(tokens) == 0 {
		return 0, fmt.Errorf("belum ada smartphone nasabah yang terdaftar di sistem (FCM Token kosong). Pastikan aplikasi Flutter di HP sudah login")
	}

	accessToken, err := s.getAccessToken()
	if err != nil {
		log.Printf("FCM access token error: %v", err)
		return 0, fmt.Errorf("gagal otentikasi Firebase FCM: %w", err)
	}

	sentCount := 0
	var lastErr error
	for _, token := range tokens {
		if err := s.sendFCM(accessToken, token, title, body); err == nil {
			sentCount++
		} else {
			log.Printf("FCM broadcast kirim gagal: %v", err)
			lastErr = err
		}
	}

	if sentCount == 0 && lastErr != nil {
		return 0, fmt.Errorf("semua token gagal dikirim ke FCM: %w", lastErr)
	}

	return sentCount, nil
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
