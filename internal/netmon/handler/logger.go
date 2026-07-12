package handler

import (
	"log"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

// LogHandler — BLOCK event'ını yerel log'a yazar.
type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) Handle(event model.BlockEvent) error {
	log.Printf("[BLOCK] PID=%d UID=%d COMM=%s -> %s:%d (%s) reason=%s policy_version=%d",
		event.Attempt.PID,
		event.Attempt.UID,
		event.Attempt.Comm,
		event.Attempt.DestIP,
		event.Attempt.DstPort,
		event.Attempt.Proto,
		event.Reason,
		event.PolicyVersion,
	)
	return nil
}
