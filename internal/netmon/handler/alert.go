package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

type alertPayload struct {
	PID     uint32 `json:"pid"`
	UID     uint32 `json:"uid"`
	Comm    string `json:"comm"`
	DestIP  string `json:"destIp"`
	DstPort uint16 `json:"dstPort"`
	Proto   string `json:"proto"`
}

type AlertHandler struct {
	client  *http.Client
	baseURL string
}

func NewAlertHandler(client *http.Client, baseURL string) *AlertHandler {
	return &AlertHandler{client: client, baseURL: baseURL}
}

func (h *AlertHandler) Handle(conn model.Connection) error {
	payload := alertPayload{
		PID:     conn.PID,
		UID:     conn.UID,
		Comm:    conn.Comm,
		DestIP:  conn.DestIP.String(),
		DstPort: conn.DstPort,
		Proto:   conn.Proto.String(),
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
