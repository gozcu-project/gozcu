package parser

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

const minSampleSize = 27

type EBPFParser struct{}

func New() *EBPFParser {
	return &EBPFParser{}
}

func (p *EBPFParser) Parse(raw []byte) (model.Connection, error) {
	if len(raw) < minSampleSize {
		return model.Connection{}, fmt.Errorf(
			"yetersiz sample boyutu: %d < %d", len(raw), minSampleSize,
		)
	}

	daddr := binary.LittleEndian.Uint32(raw[8:12])
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, daddr)

	return model.Connection{
		PID:     binary.LittleEndian.Uint32(raw[0:4]),
		UID:     binary.LittleEndian.Uint32(raw[4:8]),
		DestIP:  net.IP(b),
		DstPort: binary.LittleEndian.Uint16(raw[12:14]),
		Proto:   model.Protocol(raw[14]),
		Comm:    parseComm(raw[15:]),
	}, nil
}

func parseComm(b []byte) string {
	n := 0
	for n < len(b) && b[n] != 0 {
		n++
	}
	return string(b[:n])
}
