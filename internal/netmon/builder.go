package netmon

import (
	"net/http"
	"time"

	"github.com/gozcu-project/gozcu/internal/netmon/bpf"
	"github.com/gozcu-project/gozcu/internal/netmon/controlplane"
	"github.com/gozcu-project/gozcu/internal/netmon/dispatcher"
	"github.com/gozcu-project/gozcu/internal/netmon/handler"
	"github.com/gozcu-project/gozcu/internal/netmon/parser"
	"github.com/gozcu-project/gozcu/internal/netmon/policy"
	"github.com/gozcu-project/gozcu/internal/netmon/repository"
	"github.com/gozcu-project/gozcu/internal/netmon/store"
)

const (
	defaultCgroupPath    = "/sys/fs/cgroup"
	defaultSyncInterval  = 30 * time.Second
)

// Builder — netmon bileşenlerini bir araya getirir.
// Dependency wiring burada yapılır, business logic içermez.
type Builder struct {
	client       *http.Client
	baseURL      string
	cgroupPath   string
	syncInterval time.Duration
}

func NewBuilder(client *http.Client, baseURL string) *Builder {
	return &Builder{
		client:       client,
		baseURL:      baseURL,
		cgroupPath:   defaultCgroupPath,
		syncInterval: defaultSyncInterval,
	}
}

func (b *Builder) WithCgroupPath(path string) *Builder {
	b.cgroupPath = path
	return b
}

func (b *Builder) WithSyncInterval(d time.Duration) *Builder {
	b.syncInterval = d
	return b
}

func (b *Builder) Build() (*Runner, error) {
	// BPF yükle
	objs, err := bpf.Load()
	if err != nil {
		return nil, err
	}

	// cgroup hook'larını attach et
	links, err := objs.AttachCgroup(b.cgroupPath)
	if err != nil {
		objs.Close()
		return nil, err
	}

	// Policy store
	policyStore := store.NewBPFPolicyStore(objs.Maps())

	// Control plane — policy sync
	repo := repository.NewWhitelistRepository(b.client, b.baseURL)
	validator := policy.NewDefaultValidator()
	tracker := controlplane.NewVersionTracker()
	syncer := controlplane.NewPolicySyncer(repo, policyStore, validator, tracker)

	// İlk sync — hata olursa boş policy ile devam (fail-closed kernel'de)
	_ = syncer.Sync()
	syncer.StartSync(b.syncInterval)

	// Handler zinciri
	chain := handler.NewChain(
		handler.NewLogHandler(),
		handler.NewAlertHandler(b.client, b.baseURL),
	)

	// Link'leri Runner'a uygun tipe çevir
	runnerLinks := make([]interface{ Close() error }, len(links))
	for i, l := range links {
		runnerLinks[i] = l
	}

	return &Runner{
		objs:       objs,
		links:      runnerLinks,
		dispatcher: dispatcher.New(chain),
		parser:     parser.New(),
	}, nil
}
