package handler

import (
	"log"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
	"github.com/gozcu-project/gozcu/internal/netmon/policy"
)

type LogHandler struct {
	decision policy.Decision
}

func NewLogHandler(decision policy.Decision) *LogHandler {
	return &LogHandler{decision: decision}
}

func (h *LogHandler) Handle(conn model.Connection) error {
	log.Printf("[%s] PID=%d UID=%d COMM=%s -> %s:%d (%s)",
		h.decision, conn.PID, conn.UID, conn.Comm,
		conn.DestIP, conn.DstPort, conn.Proto)
	return nil
}
