package couchdb

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	responseRoot, _        = ioutil.ReadFile("testdata/root.json")
	responseNodeStats, _   = ioutil.ReadFile("testdata/node_stats.json")
	responseActiveTasks, _ = ioutil.ReadFile("testdata/active_tasks.json")
	responseNodeSystem, _  = ioutil.ReadFile("testdata/node_system.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"responseRoot":        responseRoot,
		"responseNodeStats":   responseNodeStats,
		"responseActiveTasks": responseActiveTasks,
		"responseNodeSystem":  responseNodeSystem,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestCouchDB_Init(t *testing.T) {
	tests := map[string]struct {
		config          Config
		wantNumOfCharts int
		wantFail        bool
	}{
		"default": {
			wantNumOfCharts: numOfCharts(
				dbActivityCharts,
				httpTrafficBreakdownCharts,
				serverOperationsCharts,
				erlangStatisticsCharts,
			),
			config: New().Config,
		},
		"URL not set": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: ""},
				}},
		},
		"invalid TLSCA": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Client: web.Client{
						TLSConfig: tlscfg.TLSConfig{TLSCA: "testdata/tls"},
					},
				}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es := New()
			es.Config = test.config

			if test.wantFail {
				assert.False(t, es.Init())
			} else {
				assert.True(t, es.Init())
				assert.Equal(t, test.wantNumOfCharts, len(*es.Charts()))
			}
		})
	}
}

func TestCouchDB_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(*testing.T) (cdb *CouchDB, cleanup func())
		wantFail bool
	}{
		"valid data":         {prepare: prepareCouchDBValidData},
		"invalid data":       {prepare: prepareCouchDBInvalidData, wantFail: true},
		"404":                {prepare: prepareCouchDB404, wantFail: true},
		"connection refused": {prepare: prepareCouchDBConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cdb, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, cdb.Check())
			} else {
				assert.True(t, cdb.Check())
			}
		})
	}
}

func numOfCharts(charts ...Charts) (num int) {
	for _, v := range charts {
		num += len(v)
	}
	return num
}

func TestCouchDB_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func prepareCouchDB(t *testing.T, createCDB func() *CouchDB) (cdb *CouchDB, cleanup func()) {
	t.Helper()
	srv := prepareCouchDBEndpoint()

	cdb = createCDB()
	cdb.URL = srv.URL
	require.True(t, cdb.Init())

	return cdb, srv.Close
}

func prepareCouchDBValidData(t *testing.T) (cdb *CouchDB, cleanup func()) {
	return prepareCouchDB(t, New)
}

func prepareCouchDBInvalidData(t *testing.T) (*CouchDB, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	cdb := New()
	cdb.URL = srv.URL
	require.True(t, cdb.Init())

	return cdb, srv.Close
}

func prepareCouchDB404(t *testing.T) (*CouchDB, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	cdb := New()
	cdb.URL = srv.URL
	require.True(t, cdb.Init())

	return cdb, srv.Close
}

func prepareCouchDBConnectionRefused(t *testing.T) (*CouchDB, func()) {
	t.Helper()
	cdb := New()
	cdb.URL = "http://127.0.0.1:38001"
	require.True(t, cdb.Init())

	return cdb, func() {}
}

func prepareCouchDBEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathOverviewStats:
				_, _ = w.Write(responseNodeStats)
			case urlPathSystemStats:
				_, _ = w.Write(responseNodeSystem)
			case urlPathActiveTasks:
				_, _ = w.Write(responseActiveTasks)
			case "/":
				_, _ = w.Write(responseRoot)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

func numOfCharts(charts ...Charts) (num int) {
	for _, v := range charts {
		num += len(v)
	}
	return num
}
