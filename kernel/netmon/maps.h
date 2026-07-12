#pragma once

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include "types.h"

/*
 * A/B double-buffer whitelist map'leri.
 * Kernel aktif olan map'e bakar, Go tarafı pasif olanı
 * doldurduktan sonra active_map pointer'ını atomik olarak değiştirir.
 */
struct {
    __uint(type, BPF_MAP_TYPE_LPM_TRIE);
    __uint(max_entries, 1024);
    __type(key, struct whitelist_key);
    __type(value, struct whitelist_value);
    __uint(map_flags, BPF_F_NO_PREALLOC);
} whitelist_a SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_LPM_TRIE);
    __uint(max_entries, 1024);
    __type(key, struct whitelist_key);
    __type(value, struct whitelist_value);
    __uint(map_flags, BPF_F_NO_PREALLOC);
} whitelist_b SEC(".maps");

/*
 * active_map — hangi whitelist map'inin aktif olduğunu tutar.
 * 0 = whitelist_a, 1 = whitelist_b
 * Go tarafı bu değeri atomik olarak güncelleyerek policy swap yapar.
 */
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u32);
} active_map SEC(".maps");

/*
 * policy_version — aktif policy'nin versiyonu.
 * BlockEvent içine eklenerek audit sırasında
 * hangi policy'nin bağlantıyı engellediği bilinir.
 */
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u32);
} policy_version SEC(".maps");

/*
 * block_events — kernel'den user-space'e BLOCK event aktarımı.
 * Ringbuf dolması enforcement'ı etkilemez:
 * kernel BLOCK kararını event publish başarısından bağımsız verir.
 */
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24); /* 16MB */
} block_events SEC(".maps");

/*
 * drop_counter — ringbuf dolu olduğunda kaçırılan event sayısı.
 * Telemetry için kullanılır, enforcement'ı etkilemez.
 */
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u64);
} drop_counter SEC(".maps");
