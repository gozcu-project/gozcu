package bpf

import (
	"fmt"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// Objects — yüklenmiş BPF nesnelerinin dışa açık temsili.
type Objects struct {
	inner *NetmonObjects
}

// Maps — kernel map referanslarını dışa açar.
type Maps struct {
	WhitelistA    *ebpf.Map
	WhitelistB    *ebpf.Map
	ActiveMap     *ebpf.Map
	PolicyVersion *ebpf.Map
	BlockEvents   *ebpf.Map
	DropCounter   *ebpf.Map
}

// Load — BPF programlarını kernel'e yükler.
func Load() (*Objects, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, fmt.Errorf("memlock kaldırılamadı: %w", err)
	}

	objs := &NetmonObjects{}
	if err := LoadNetmonObjects(objs, nil); err != nil {
		return nil, fmt.Errorf("BPF nesneleri yüklenemedi: %w", err)
	}

	return &Objects{inner: objs}, nil
}

// Close — kernel kaynaklarını serbest bırakır.
func (o *Objects) Close() {
	o.inner.Close()
}

// Maps — kernel map referanslarını döner.
func (o *Objects) Maps() Maps {
	return Maps{
		WhitelistA:    o.inner.WhitelistA,
		WhitelistB:    o.inner.WhitelistB,
		ActiveMap:     o.inner.ActiveMap,
		PolicyVersion: o.inner.PolicyVersion,
		BlockEvents:   o.inner.BlockEvents,
		DropCounter:   o.inner.DropCounter,
	}
}

// AttachCgroup — cgroup/connect4 ve cgroup/connect6 hook'larını attach eder.
func (o *Objects) AttachCgroup(cgroupPath string) ([]link.Link, error) {
	if err := validateCgroup(cgroupPath); err != nil {
		return nil, err
	}

	connect4, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInet4Connect,
		Program: o.inner.GozcuConnect4,
	})
	if err != nil {
		return nil, fmt.Errorf("cgroup/connect4 attach edilemedi: %w", err)
	}

	connect6, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInet6Connect,
		Program: o.inner.GozcuConnect6,
	})
	if err != nil {
		connect4.Close()
		return nil, fmt.Errorf("cgroup/connect6 attach edilemedi: %w", err)
	}

	return []link.Link{connect4, connect6}, nil
}

func validateCgroup(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("cgroup path erişilemiyor %s: %w", path, err)
	}
	return nil
}
