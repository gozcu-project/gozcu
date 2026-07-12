package model

// BlockEvent — kernel'den ringbuf üzerinden iletilen
// engelleme event'ının domain temsili.
//
// Tasarım notu: ALLOW event'ları iletilmez.
// Kernel yalnızca BLOCK kararlarını ringbuf'a yazar,
// bu sayede yüksek trafikli ortamlarda overhead minimize edilir.
type BlockEvent struct {
	Attempt       ConnectionAttempt
	Reason        BlockReason
	PolicyVersion uint32
}
