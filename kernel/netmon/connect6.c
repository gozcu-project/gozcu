//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>
#include "types.h"
#include "maps.h"
#include "events.h"

#define AF_INET6 10

static __always_inline __u32 ipv6_to_ipv4(const __u32 ip6[4])
{
    if (ip6[0] == 0 && ip6[1] == 0 && ip6[2] == bpf_htonl(0x0000ffff))
        return ip6[3];
    return 0;
}

static __always_inline void fill_attempt6(
    struct bpf_sock_addr      *ctx,
    struct connection_attempt *out
) {
    out->pid      = bpf_get_current_pid_tgid() >> 32;
    out->uid      = bpf_get_current_uid_gid() & 0xFFFFFFFF;
    out->dst_ip   = ipv6_to_ipv4(ctx->user_ip6);
    out->dst_port = bpf_ntohs(ctx->user_port >> 16);
    out->proto    = ctx->protocol;
    bpf_get_current_comm(&out->comm, sizeof(out->comm));
}

static __always_inline int evaluate6(
    const struct connection_attempt *attempt
) {
    if (attempt->dst_ip == 0)
        return BLOCK_REASON_POLICY_ERROR;

    const struct whitelist_value *val = lookup_whitelist(attempt->dst_ip);

    if (!val)
        return BLOCK_REASON_WHITELIST_MISS;

    return port_allowed(val, attempt->dst_port)
        ? 0
        : BLOCK_REASON_WHITELIST_MISS;
}

SEC("cgroup/connect6")
int gozcu_connect6(struct bpf_sock_addr *ctx)
{
    if (ctx->family != AF_INET6)
        return 1;

    struct connection_attempt attempt = {};
    fill_attempt6(ctx, &attempt);

    return enforce(&attempt, evaluate6(&attempt));
}
