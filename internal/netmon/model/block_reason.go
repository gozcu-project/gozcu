package model

// BlockReason — bağlantının neden engellendiğini açıklar.
// Kernel'deki enum block_reason ile birebir eşleşmelidir.
type BlockReason uint32

const (
	BlockReasonWhitelistMiss BlockReason = 0
	BlockReasonPolicyError   BlockReason = 1
	BlockReasonNoActiveMap   BlockReason = 2
)

func (r BlockReason) String() string {
	names := map[BlockReason]string{
		BlockReasonWhitelistMiss: "WHITELIST_MISS",
		BlockReasonPolicyError:   "POLICY_ERROR",
		BlockReasonNoActiveMap:   "NO_ACTIVE_MAP",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return "UNKNOWN"
}
