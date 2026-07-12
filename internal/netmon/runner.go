package netmon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/ringbuf"

	"github.com/gozcu-project/gozcu/internal/netmon/bpf"
	"github.com/gozcu-project/gozcu/internal/netmon/dispatcher"
	"github.com/gozcu-project/gozcu/internal/netmon/parser"
)

// Runner — ringbuf event döngüsünü yöneten bileşen.
// Tek sorumluluğu: ringbuf'tan raw event oku → parse et → dispatch et.
// Policy sync, map management, alert retry Runner'ın sorumluluğunda değil.
type Runner struct {
	objs       *bpf.Objects
	links      []interface{ Close() error }
	dispatcher *dispatcher.Dispatcher
	parser     *parser.BlockEventParser
}

// Run — event döngüsünü başlatır, SIGINT/SIGTERM ile durur.
func (r *Runner) Run() error {
	maps := r.objs.Maps()

	rd, err := ringbuf.NewReader(maps.BlockEvents)
	if err != nil {
		return err
	}
	defer rd.Close()

	log.Println("[runner] gozcu-netmon başladı, enforcement aktif")

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
			log.Printf("[runner] ringbuf okuma hatası: %v", err)
			continue
		}

		event, err := r.parser.Parse(record.RawSample)
		if err != nil {
			log.Printf("[runner] event parse hatası: %v", err)
			continue
		}

		r.dispatcher.Dispatch(event)
	}
}

// Close — BPF kaynaklarını serbest bırakır.
func (r *Runner) Close() {
	for _, l := range r.links {
		l.Close()
	}
	r.objs.Close()
}
