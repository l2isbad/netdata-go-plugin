package couchdb

import (
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("couchdb", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 1,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *CouchDB {
	return &CouchDB{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{
					URL: "http://127.0.0.1:5984",
				},
				Client: web.Client{
					Timeout: web.Duration{Duration: time.Second * 5},
				},
			},
		},
		collectedIndices: make(map[string]bool),
	}
}

type (
	Config struct {
		web.HTTP `yaml:",inline"`
	}

	CouchDB struct {
		module.Base
		Config `yaml:",inline"`

		httpClient       *http.Client
		charts           *module.Charts
		collectedIndices map[string]bool
	}
)

func (cdb *CouchDB) Cleanup() {
	if cdb.httpClient == nil {
		return
	}
	cdb.httpClient.CloseIdleConnections()
}

func (cdb *CouchDB) Init() bool {
	err := cdb.validateConfig()
	if err != nil {
		cdb.Errorf("check configuration: %v", err)
		return false
	}

	httpClient, err := cdb.initHTTPClient()
	if err != nil {
		cdb.Errorf("init HTTP client: %v", err)
		return false
	}
	cdb.httpClient = httpClient

	charts, err := cdb.initCharts()
	if err != nil {
		cdb.Errorf("init charts: %v", err)
		return false
	}
	cdb.charts = charts

	return true
}

func (cdb *CouchDB) Check() bool {
	if err := cdb.pingCouchDB(); err != nil {
		cdb.Error(err)
		return false
	}
	return true //TODO: len(cdb.Collect()) > 0
}

func (cdb *CouchDB) Charts() *Charts {
	return cdb.charts
}

func (cdb *CouchDB) Collect() map[string]int64 {
	mx, err := cdb.collect()
	if err != nil {
		cdb.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}
