package parser

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

/*
 * Kernel struct block_event layout (types.h ile birebir eşleşmeli):
 *
 * struct connection_attempt {
 *     __u32 pid;       offset 0
 *     __u32 uid;       offset 4
 *     __u32 dst_ip;    offset 8
 *     __u16 dst_port;  offset 12
 *     __u8  proto;     offset 14
 *     char  comm[16];  offset 15
 * };                   size: 31 bytes (+ 1 byte padding = 32)
 *
 * struct block_event {
 *     struct connection_attempt attempt;  offset 0
 *     enum block_reason reason;           offset 32
 *     __u32 policy_version;              offset 36
 * };                                      size: 40 bytes
 */
const (
	offsetPID           = 0
	offsetUID           = 4
	offsetDstIP         = 8
	offsetDstPort       = 12
	offsetProto         = 14
	offsetComm          = 15
	offsetReason        = 32
	offsetPolicyVersion = 36
	minEventSize        = 40
)

// BlockEventParser — ham BPF ringbuf sample'ını BlockEvent'a dönüştürür.
// Kernel struct layout'u ile birebir senkronize olmalıdır.
type BlockEventParser struct{}

func New() *BlockEventParser {
	return &BlockEventParser{}
}

// Parse — ham byte slice'ını BlockEvent'a dönüştürür.
// Hatalı boyut veya geçersiz veri durumunda hata döner.
func (p *BlockEventParser) Parse(raw []byte) (model.BlockEvent, error) {
	if len(raw) < minEventSize {
		return model.BlockEvent{}, fmt.Errorf(
			"yetersiz event boyutu: %d < %d", len(raw), minEventSize,
		)
	}

	return model.BlockEvent{
		Attempt: model.ConnectionAttempt{
			PID:     binary.LittleEndian.Uint32(raw[offsetPID:]),
			UID:     binary.LittleEndian.Uint32(raw[offsetUID:]),
			DestIP:  parseIP(raw[offsetDstIP:]),
			DstPort: binary.LittleEndian.Uint16(raw[offsetDstPort:]),
			Proto:   model.Protocol(raw[offsetProto]),
			Comm:    parseComm(raw[offsetComm:]),
		},
		Reason:        model.BlockReason(binary.LittleEndian.Uint32(raw[offsetReason:])),
		PolicyVersion: binary.LittleEndian.Uint32(raw[offsetPolicyVersion:]),
	}, nil
}

func parseIP(b []byte) net.IP {
	ip := make([]byte, 4)
	binary.LittleEndian.PutUint32(ip, binary.LittleEndian.Uint32(b))
	return net.IP(ip)
}

func parseComm(b []byte) string {
	n := 0
	for n < len(b) && n < 16 && b[n] != 0 {
		n++
	}
	return string(b[:n])
}
