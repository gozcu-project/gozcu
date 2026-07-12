package model

import "net"

// ConnectionAttempt — cgroup/connect hook'unda yakalanan
// bağlantı girişiminin değişmez domain temsili.
// "Connection" değil "Attempt" — başarılı bir bağlantı değil,
// kernel tarafından yakalanan bir girişimdir.
type ConnectionAttempt struct {
	PID     uint32
	UID     uint32
	Comm    string
	DestIP  net.IP
	DstPort uint16
	Proto   Protocol
}
