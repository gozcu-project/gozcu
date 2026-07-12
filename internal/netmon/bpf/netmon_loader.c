//go:build ignore

/*
 * netmon_loader.c — bpf2go için tek entry point.
 * Tüm kernel eBPF kaynak dosyaları buradan include edilir.
 * Bu dosya hiçbir zaman manuel düzenlenmez.
 */
#include "../../../kernel/netmon/connect4.c"
#include "../../../kernel/netmon/connect6.c"

char LICENSE[] SEC("license") = "GPL";
