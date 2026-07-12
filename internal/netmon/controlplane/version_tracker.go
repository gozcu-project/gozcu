package controlplane

import "sync/atomic"

// VersionTracker — aktif policy version'ını thread-safe takip eder.
// Agent eski version'a geri dönemez — rollback sadece
// control plane tarafından yapılabilir (version artırılarak).
type VersionTracker struct {
	current atomic.Uint32
}

func NewVersionTracker() *VersionTracker {
	return &VersionTracker{}
}

// IsCurrent — verilen version zaten aktifse true döner.
func (t *VersionTracker) IsCurrent(version uint32) bool {
	return t.current.Load() == version
}

// Update — yeni version'ı kaydeder.
func (t *VersionTracker) Update(version uint32) {
	t.current.Store(version)
}

// Current — aktif version'ı döner.
func (t *VersionTracker) Current() uint32 {
	return t.current.Load()
}
