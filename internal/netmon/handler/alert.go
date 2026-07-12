package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

type alertPayload struct {
	PID           uint32 `json:"pid"`
	UID           uint32 `json:"uid"`
	Comm          string `json:"comm"`
	DestIP        string `json:"destIp"`
	DstPort       uint16 `json:"dstPort"`
	Proto         string `json:"proto"`
	Reason        string `json:"reason"`
	PolicyVersion uint32 `json:"policyVersion"`
}

// AlertHandler — BLOCK event'ını mTLS üzerinden Spring Boot'a bildirir.
// Gönderme başarısızlığı enforcement'ı etkilemez —
// hata döner, Chain devam eder.
type AlertHandler struct {
	client  *http.Client
	baseURL string
}

func NewAlertHandler(client *http.Client, baseURL string) *AlertHandler {
	return &AlertHandler{client: client, baseURL: baseURL}
}

func (h *AlertHandler) Handle(event model.BlockEvent) error {
	payload := alertPayload{
		PID:           event.Attempt.PID,
		UID:           event.Attempt.UID,
		Comm:          event.Attempt.Comm,
		DestIP:        event.Attempt.DestIP.String(),
		DstPort:       event.Attempt.DstPort,
		Proto:         event.Attempt.Proto.String(),
		Reason:        event.Reason.String(),
		PolicyVersion: event.PolicyVersion,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alarm serialize edilemedi: %w", err)
	}

	resp, err := h.client.Post(
		fmt.Sprintf("%s/api/netmon/alerts", h.baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("alarm gönderilemedi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend alarm için %d döndü", resp.StatusCode)
	}
	return nil
}
