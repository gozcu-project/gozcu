package policy

// Policy — whitelist kurallarının domain temsili.
// Go tarafı ALLOW/BLOCK kararı vermez — bu kernel'in işi.
// Policy burada yalnızca BPF map'e yüklenecek kural setini tanımlar.
type Policy struct {
	Version uint32
	Rules   []Rule
}
