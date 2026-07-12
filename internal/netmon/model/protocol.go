package model

// Protocol — tip-güvenli ağ protokolü tanımı.
type Protocol uint8

const (
	ProtoTCP Protocol = 6
	ProtoUDP Protocol = 17
)

func (p Protocol) String() string {
	names := map[Protocol]string{
		ProtoTCP: "TCP",
		ProtoUDP: "UDP",
	}
	if name, ok := names[p]; ok {
		return name
	}
	return "UNKNOWN"
}
