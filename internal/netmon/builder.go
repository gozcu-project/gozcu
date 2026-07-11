package netmon

import (
	"log"
	"net/http"
	"time"

	"github.com/gozcu-project/gozcu/internal/netmon/dispatcher"
	"github.com/gozcu-project/gozcu/internal/netmon/handler"
	"github.com/gozcu-project/gozcu/internal/netmon/parser"
	"github.com/gozcu-project/gozcu/internal/netmon/policy"
	"github.com/gozcu-project/gozcu/internal/netmon/repository"
)

type Builder struct {
	client          *http.Client
	baseURL         string
	refreshInterval time.Duration
}

func NewBuilder(client *http.Client, baseURL string) *Builder {
	return &Builder{
		client:          client,
		baseURL:         baseURL,
		refreshInterval: 30 * time.Second,
	}
}

func (b *Builder) WithRefreshInterval(d time.Duration) *Builder {
	b.refreshInterval = d
	return b
}

func (b *Builder) Build() (*Runner, error) {
	wp := policy.NewWhitelistPolicy([]policy.Rule{})

	repo := repository.NewWhitelistRepository(b.client, b.baseURL)
	if rules, err := repo.Load(); err != nil {
		log.Printf("whitelist henüz hazır değil, boş policy ile devam")
	} else {
		wp.Update(rules)
	}
	repo.StartSync(wp, b.refreshInterval)

	d := dispatcher.New().
		WithHandler(policy.Allow, handler.NewChain(
			handler.NewLogHandler(policy.Allow),
		)).
		WithHandler(policy.Block, handler.NewChain(
			handler.NewLogHandler(policy.Block),
			// AlertHandler backend endpoint'leri hazır olunca eklenecek
		))

	return &Runner{
		policy:     wp,
		dispatcher: d,
		parser:     parser.New(),
	}, nil
}
