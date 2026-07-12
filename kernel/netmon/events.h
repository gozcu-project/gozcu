#pragma once

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include "types.h"
#include "maps.h"

static __always_inline void emit_block_event(
    const struct connection_attempt *attempt,
    enum block_reason reason
) {
    __u32 key = 0;
    __u32 *ver = bpf_map_lookup_elem(&policy_version, &key);
    __u32 pver = ver ? *ver : 0;

    struct block_event *e = bpf_ringbuf_reserve(&block_events, sizeof(*e), 0);
    if (!e) {
        __u64 *cnt = bpf_map_lookup_elem(&drop_counter, &key);
        if (cnt)
            __sync_fetch_and_add(cnt, 1);
        return;
    }

    e->attempt        = *attempt;
    e->reason         = reason;
    e->policy_version = pver;

    bpf_ringbuf_submit(e, 0);
}

static __always_inline const struct whitelist_value *lookup_whitelist(
    __u32 dst_ip
) {
    __u32 key = 0;
    __u32 *active = bpf_map_lookup_elem(&active_map, &key);
    if (!active)
        return NULL;

    struct whitelist_key wkey = {
        .prefix_len = 32,
        .ip         = dst_ip,
    };

    if (*active == 0)
        return bpf_map_lookup_elem(&whitelist_a, &wkey);
    else
        return bpf_map_lookup_elem(&whitelist_b, &wkey);
}

static __always_inline int port_allowed(
    const struct whitelist_value *val,
    __u16 dst_port
) {
    if (val->port_count == 0)
        return 1;

    #pragma unroll
    for (int i = 0; i < MAX_PORTS; i++) {
        if (i >= val->port_count)
            break;
        if (val->ports[i] == dst_port)
            return 1;
    }
    return 0;
}

/*
 * enforce — paylaşılan karar uygulama fonksiyonu.
 * connect4 ve connect6 tarafından kullanılır.
 */
static __always_inline int enforce(
    const struct connection_attempt *attempt,
    int reason
) {
    if (reason == 0)
        return 1;

    emit_block_event(attempt, (enum block_reason)reason);
    return -1; /* -EPERM */
}
