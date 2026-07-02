//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>
#include <bpf/bpf_tracing.h>

#define AF_INET 2

struct event {
    __u32 pid;
    __u32 uid;
    __u32 daddr;
    __u16 dport;
    __u8  proto;
    char  comm[16];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
} events SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_connect")
int trace_connect(struct trace_event_raw_sys_enter *ctx)
{
    struct sockaddr_in sa = {};

    // user-space pointer'dan güvenli okuma
    void *uservaddr = (void *)ctx->args[1];
    if (bpf_probe_read_user(&sa, sizeof(sa), uservaddr) != 0)
        return 0;

    if (sa.sin_family != AF_INET)
        return 0;

    struct event *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e)
        return 0;

    __u64 pid_tgid = bpf_get_current_pid_tgid();
    e->pid   = pid_tgid >> 32;
    e->uid   = bpf_get_current_uid_gid() & 0xFFFFFFFF;
    e->daddr = sa.sin_addr.s_addr;
    e->dport = bpf_ntohs(sa.sin_port);
    e->proto = 6;

    bpf_get_current_comm(&e->comm, sizeof(e->comm));

    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
