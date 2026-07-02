package bpf

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

type Event struct {
	PID   uint32
	UID   uint32
	DAddr uint32
	DPort uint16
	Proto uint8
	Comm  [16]byte
}

func (e *Event) DestIP() string {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, e.DAddr)
	return net.IP(b).String()
}

func (e *Event) CommStr() string {
	n := 0
	for n < len(e.Comm) && e.Comm[n] != 0 {
		n++
	}
	return string(e.Comm[:n])
}

func Run(onEvent func(e *Event)) error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("memlock kaldırılamadı: %w", err)
	}

	objs := netmonObjects{}
	if err := loadNetmonObjects(&objs, nil); err != nil {
		return fmt.Errorf("eBPF objeleri yüklenemedi: %w", err)
	}
	defer objs.Close()

	tp, err := link.Tracepoint("syscalls", "sys_enter_connect", objs.TraceConnect, nil)
	if err != nil {
		return fmt.Errorf("tracepoint bağlanamadı: %w", err)
	}
	defer tp.Close()

	rd, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		return fmt.Errorf("ring buffer açılamadı: %w", err)
	}
	defer rd.Close()

	log.Println("gozcu-netmon başladı, network olayları izleniyor...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		rd.Close()
	}()

	for {
		record, err := rd.Read()
		if err != nil {
			if err == ringbuf.ErrClosed {
				return nil
			}
			log.Printf("ring buffer okuma hatası: %v", err)
			continue
		}


		if len(record.RawSample) < 27 {
			continue
		}

		e := &Event{}
		e.PID   = binary.LittleEndian.Uint32(record.RawSample[0:4])
		e.UID   = binary.LittleEndian.Uint32(record.RawSample[4:8])
		e.DAddr = binary.LittleEndian.Uint32(record.RawSample[8:12])
		e.DPort = binary.LittleEndian.Uint16(record.RawSample[12:14])
		e.Proto = record.RawSample[14]
		copy(e.Comm[:], record.RawSample[15:])

		if onEvent != nil {
			onEvent(e)
		}
	}
}
