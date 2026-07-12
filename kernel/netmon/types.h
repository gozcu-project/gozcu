#pragma once

#include "vmlinux.h"

/*
 * WhitelistKey — BPF_MAP_TYPE_LPM_TRIE için arama anahtarı.
 * LPM_TRIE prefix uzunluğu + data alanı bekler.
 * prefix_len: anlamlı bit sayısı (örn. 24 → /24 ağ maskesi)
 * ip:         IPv4 adresi (network byte order)
 */
struct whitelist_key {
    __u32 prefix_len;
    __u32 ip;
};

/*
 * WhitelistValue — eşleşen kural için port bilgisi.
 * port_count = 0 ise tüm portlara izin verilir.
 * Maksimum 16 port tanımlanabilir.
 */
#define MAX_PORTS 16

struct whitelist_value {
    __u16 ports[MAX_PORTS];
    __u8  port_count;
};

/*
 * BlockReason — bağlantı neden engellendi?
 */
enum block_reason {
    BLOCK_REASON_WHITELIST_MISS  = 0,
    BLOCK_REASON_POLICY_ERROR    = 1,
    BLOCK_REASON_NO_ACTIVE_MAP   = 2,
};

/*
 * ConnectionAttempt — bir bağlantı girişiminin kernel tarafında
 * yakalandığı andaki değişmez temsili.
 */
struct connection_attempt {
    __u32 pid;
    __u32 uid;
    __u32 dst_ip;
    __u16 dst_port;
    __u8  proto;
    char  comm[16];
};

/*
 * BlockEvent — kernel'den user-space'e ringbuf üzerinden
 * iletilen BLOCK event'ı. ALLOW event'ları iletilmez.
 */
struct block_event {
    struct connection_attempt attempt;
    enum block_reason         reason;
    __u32                     policy_version;
};
