package coredns

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultURL         = "http://127.0.0.1:9153/metrics"
	defaultHTTPTimeout = time.Second * 2
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("coredns", creator)
}

// New creates CoreDNS with default values.
func New() *CoreDNS {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		//PerZoneStats:   filter{Include: []string{"glob:*"}},
		//PerServerStats: filter{Include: []string{"glob:*"}},
	}
	return &CoreDNS{
		Config:           config,
		charts:           summaryCharts.Copy(),
		collectedServers: make(map[string]bool),
		collectedZones:   make(map[string]bool),
	}
}

// Config is the CoreDNS module configuration.
type Config struct {
	web.HTTP       `yaml:",inline"`
	PerServerStats filter `yaml:"per_server_stats"`
	PerZoneStats   filter `yaml:"per_zone_stats"`
}

// CoreDNS CoreDNS module.
type CoreDNS struct {
	module.Base
	Config           `yaml:",inline"`
	charts           *Charts
	prom             prometheus.Prometheus
	perServerMatcher matcher.Matcher
	perZoneMatcher   matcher.Matcher
	collectedServers map[string]bool
	collectedZones   map[string]bool
}

// Cleanup makes cleanup.
func (CoreDNS) Cleanup() {}

// Init makes initialization.
func (cd *CoreDNS) Init() bool {
	if cd.URL == "" {
		cd.Error("URL parameter is not set")
		return false
	}

	if !cd.PerServerStats.isEmpty() {
		m, err := cd.PerServerStats.createMatcher()
		if err != nil {
			cd.Errorf("error on creating 'per_server_stats' matcher : %v", err)
			return false
		}
		cd.perServerMatcher = matcher.WithCache(m)
	}

	if !cd.PerZoneStats.isEmpty() {
		m, err := cd.PerZoneStats.createMatcher()
		if err != nil {
			cd.Errorf("error on creating 'per_zone_stats' matcher : %v", err)
			return false
		}
		cd.perZoneMatcher = matcher.WithCache(m)
	}

	client, err := web.NewHTTPClient(cd.Client)
	if err != nil {
		cd.Errorf("error on creating http client : %v", err)
		return false
	}

	cd.prom = prometheus.New(client, cd.Request)

	return true
}

// Check makes check.
func (cd CoreDNS) Check() bool {
	return len(cd.Collect()) > 0
}

// Charts creates Charts.
func (cd CoreDNS) Charts() *Charts {
	return cd.charts
}

// Collect collects metrics.
func (cd *CoreDNS) Collect() map[string]int64 {
	mx, err := cd.collect()

	if err != nil {
		cd.Error(err)
		return nil
	}

	return mx
}
