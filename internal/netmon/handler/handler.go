package handler

import "github.com/gozcu-project/gozcu/internal/netmon/model"

// EventHandler — bir BlockEvent'ı işleyen sözleşme.
// Chain of Responsibility: her handler kendi işini yapar,
// sonra zincirdeki bir sonraki handler'a iletir.
// Alert gönderme başarısızlığı enforcement'ı etkilemez —
// handler hataları zinciri durdurmaz, loglanır.
type EventHandler interface {
	Handle(event model.BlockEvent) error
}
