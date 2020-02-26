package vernemq

import (
	"errors"
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("vernemq", creator)
}

func New() *VerneMQ {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: "http://127.0.0.1:8888/metrics",
			},
			Client: web.Client{Timeout: web.Duration{Duration: time.Second}},
		},
	}

	return &VerneMQ{
		Config: config,
		charts: charts.Copy(),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	VerneMQ struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *Charts
	}
)

func (v VerneMQ) validateConfig() error {
	if v.UserURL == "" {
		return errors.New("URL is not set")
	}
	return nil
}

func (v *VerneMQ) initClient() error {
	client, err := web.NewHTTPClient(v.Client)
	if err != nil {
		return err
	}

	v.prom = prometheus.New(client, v.Request)
	return nil
}

func (v *VerneMQ) Init() bool {
	if err := v.validateConfig(); err != nil {
		v.Errorf("error on validating config: %v", err)
		return false
	}
	if err := v.initClient(); err != nil {
		v.Errorf("error on initializing client: %v", err)
		return false
	}
	return true
}

func (v *VerneMQ) Check() bool {
	return len(v.Collect()) > 0
}

func (v *VerneMQ) Charts() *Charts {
	return v.charts
}

func (v *VerneMQ) Collect() map[string]int64 {
	mx, err := v.collect()
	if err != nil {
		v.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (VerneMQ) Cleanup() {}
