package netmon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"

	"github.com/gozcu-project/gozcu/internal/netmon/bpf"
	"github.com/gozcu-project/gozcu/internal/netmon/dispatcher"
	"github.com/gozcu-project/gozcu/internal/netmon/parser"
	"github.com/gozcu-project/gozcu/internal/netmon/policy"
)

type Runner struct {
	policy     policy.Policy
	dispatcher *dispatcher.Dispatcher
	parser     *parser.EBPFParser
}

func (r *Runner) Run() error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return err
	}

	objs, err := bpf.LoadObjects()
	if err != nil {
		return err
	}
	defer objs.Close()

	tp, err := objs.AttachTracepoint()
	if err != nil {
		return err
	}
	defer tp.Close()

	rd, err := ringbuf.NewReader(objs.Events().Events)
	if err != nil {
		return err
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

		conn, err := r.parser.Parse(record.RawSample)
		if err != nil {
			continue
		}

		decision := r.policy.Evaluate(conn)
		_ = r.dispatcher.Dispatch(decision, conn)
	}
}
