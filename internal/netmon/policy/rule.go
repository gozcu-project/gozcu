package policy

import "net"

// Rule — tek bir whitelist kuralı.
// CIDR: hedef IP aralığı (örn. 0.0.0.0/0, 10.0.0.0/8)
// Ports: izin verilen portlar. Boşsa tüm portlara izin verilir.
type Rule struct {
	CIDR  *net.IPNet
	Ports []uint16
}
