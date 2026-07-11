package bpf

import "github.com/cilium/ebpf/link"

// Objects — dışarıya açık wrapper, skeleton'ın unexported tipini sarar.
type Objects struct {
	inner netmonObjects
}

func LoadObjects() (*Objects, error) {
	o := &Objects{}
	if err := loadNetmonObjects(&o.inner, nil); err != nil {
		return nil, err
	}
	return o, nil
}

func (o *Objects) Close() {
	o.inner.Close()
}

func (o *Objects) AttachTracepoint() (link.Link, error) {
	return link.Tracepoint("syscalls", "sys_enter_connect", o.inner.TraceConnect, nil)
}

func (o *Objects) Events() *netmonMaps {
	return &o.inner.netmonMaps
}
