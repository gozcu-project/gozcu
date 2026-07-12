//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>
#include "types.h"
#include "maps.h"
#include "events.h"

#define AF_INET 2

static __always_inline void fill_attempt(
    struct bpf_sock_addr      *ctx,
    struct connection_attempt *out
) {
    out->pid      = bpf_get_current_pid_tgid() >> 32;
    out->uid      = bpf_get_current_uid_gid() & 0xFFFFFFFF;
    out->dst_ip   = ctx->user_ip4;
    out->dst_port = bpf_ntohs(ctx->user_port >> 16);
    out->proto    = ctx->protocol;
    bpf_get_current_comm(&out->comm, sizeof(out->comm));
}

static __always_inline int evaluate(
    const struct connection_attempt *attempt
) {
    const struct whitelist_value *val = lookup_whitelist(attempt->dst_ip);

    if (!val)
        return BLOCK_REASON_WHITELIST_MISS;

    return port_allowed(val, attempt->dst_port)
        ? 0
        : BLOCK_REASON_WHITELIST_MISS;
}

SEC("cgroup/connect4")
int gozcu_connect4(struct bpf_sock_addr *ctx)
{
    if (ctx->family != AF_INET)
        return 1;

    struct connection_attempt attempt = {};
    fill_attempt(ctx, &attempt);

    return enforce(&attempt, evaluate(&attempt));
}
