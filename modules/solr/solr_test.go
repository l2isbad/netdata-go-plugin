package solr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	coreMetricsV6, _ = ioutil.ReadFile("testdata/core-metrics-v6.txt")
	coreMetricsV7, _ = ioutil.ReadFile("testdata/core-metrics-v7.txt")
)

func version(v string) string {
	return format(`{ "lucene":{ "solr-spec-version":"%s"}}`, v)
}

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*module.Module)(nil), mod)
	assert.Equal(t, defaultURL, mod.URL)
	assert.Equal(t, defaultHTTPTimeout, mod.Client.Timeout.Duration)
}

func TestSolr_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())

	assert.NotNil(t, mod.reqInfoSystem)
	assert.NotNil(t, mod.reqCoreHandlers)
	assert.NotNil(t, mod.client)
}

func TestSolr_Check(t *testing.T) {
	mod := New()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/solr/admin/info/system" {
					_, _ = w.Write([]byte(version(fmt.Sprintf("%.1f.0", minSupportedVersion))))
					return
				}
			}))

	mod.URL = ts.URL

	require.True(t, mod.Init())

	assert.True(t, mod.Check())
}

func TestSolr_Check_UnsupportedVersion(t *testing.T) {
	mod := New()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/solr/admin/info/system" {
					_, _ = w.Write([]byte(version(fmt.Sprintf("%.1f.0", minSupportedVersion-1))))
					return
				}
			}))

	mod.URL = ts.URL

	require.True(t, mod.Init())

	assert.False(t, mod.Check())
}

func TestSolr_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestSolr_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestSolr_CollectV6(t *testing.T) {
	mod := New()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/solr/admin/info/system" {
					_, _ = w.Write([]byte(version(fmt.Sprintf("%.1f.0", minSupportedVersion))))
					return
				}
				if r.URL.Path == "/solr/admin/metrics" {
					_, _ = w.Write(coreMetricsV6)
					return
				}
			}))

	mod.URL = ts.URL

	require.True(t, mod.Init())
	require.True(t, mod.Check())
	require.NotNil(t, mod.Charts())

	expected := map[string]int64{
		"core2_query_requestTimes_min_ms":     0,
		"core1_query_serverErrors_count":      3,
		"core2_update_requestTimes_mean_ms":   0,
		"core2_query_requestTimes_p99_ms":     297000000,
		"core2_query_requestTimes_p999_ms":    2997000000,
		"core1_update_requestTimes_p99_ms":    297000000,
		"core2_update_requestTimes_p75_ms":    225000000,
		"core2_update_requests_count":         3,
		"core2_query_requestTimes_p75_ms":     225000000,
		"core2_update_requestTimes_min_ms":    0,
		"core2_query_clientErrors_count":      3,
		"core2_query_requestTimes_count":      3,
		"core2_query_requestTimes_median_ms":  0,
		"core2_query_requestTimes_p95_ms":     285000000,
		"core2_update_serverErrors_count":     3,
		"core1_query_requestTimes_mean_ms":    0,
		"core1_update_totalTime_count":        3,
		"core1_update_errors_count":           3,
		"core1_query_errors_count":            3,
		"core1_query_timeouts_count":          3,
		"core1_update_requestTimes_p95_ms":    285000000,
		"core1_query_clientErrors_count":      3,
		"core2_query_serverErrors_count":      3,
		"core1_update_requestTimes_p75_ms":    225000000,
		"core2_update_requestTimes_p99_ms":    297000000,
		"core2_query_requests_count":          3,
		"core2_update_clientErrors_count":     3,
		"core1_update_requestTimes_min_ms":    0,
		"core1_update_requestTimes_mean_ms":   0,
		"core1_query_requestTimes_p95_ms":     285000000,
		"core1_query_requestTimes_p999_ms":    2997000000,
		"core1_update_serverErrors_count":     3,
		"core1_query_requests_count":          3,
		"core1_update_requestTimes_p999_ms":   2997000000,
		"core1_query_requestTimes_p75_ms":     225000000,
		"core1_update_requestTimes_count":     3,
		"core2_update_requestTimes_p95_ms":    285000000,
		"core1_query_requestTimes_count":      3,
		"core1_query_requestTimes_p99_ms":     297000000,
		"core1_update_requestTimes_median_ms": 0,
		"core1_update_requestTimes_max_ms":    0,
		"core2_update_requestTimes_count":     3,
		"core1_query_requestTimes_min_ms":     0,
		"core1_update_timeouts_count":         3,
		"core2_update_timeouts_count":         3,
		"core2_update_errors_count":           3,
		"core1_update_requests_count":         3,
		"core2_query_errors_count":            3,
		"core1_query_requestTimes_median_ms":  0,
		"core1_query_requestTimes_max_ms":     0,
		"core1_update_clientErrors_count":     3,
		"core2_update_requestTimes_median_ms": 0,
		"core2_query_requestTimes_mean_ms":    0,
		"core2_update_totalTime_count":        3,
		"core2_update_requestTimes_max_ms":    0,
		"core2_update_requestTimes_p999_ms":   2997000000,
		"core2_query_timeouts_count":          3,
		"core2_query_requestTimes_max_ms":     0,
		"core1_query_totalTime_count":         3,
		"core2_query_totalTime_count":         3,
	}

	assert.Equal(t, expected, mod.Collect())
	assert.Equal(t, expected, mod.Collect())
}

func TestSolr_CollectV7(t *testing.T) {
	mod := New()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/solr/admin/info/system" {
					_, _ = w.Write([]byte(version(fmt.Sprintf("%.1f.0", minSupportedVersion+1))))
					return
				}
				if r.URL.Path == "/solr/admin/metrics" {
					_, _ = w.Write(coreMetricsV7)
					return
				}
			}))

	mod.URL = ts.URL

	require.True(t, mod.Init())
	require.True(t, mod.Check())
	require.NotNil(t, mod.Charts())

	expected := map[string]int64{
		"core1_query_requestTimes_p95_ms":     285000000,
		"core1_query_timeouts_count":          3,
		"core1_update_requestTimes_p999_ms":   2997000000,
		"core2_query_requestTimes_mean_ms":    0,
		"core2_query_timeouts_count":          3,
		"core1_update_timeouts_count":         3,
		"core1_update_requestTimes_mean_ms":   0,
		"core2_update_serverErrors_count":     3,
		"core2_query_requestTimes_min_ms":     0,
		"core2_query_requestTimes_p75_ms":     225000000,
		"core2_update_clientErrors_count":     3,
		"core2_update_requestTimes_count":     3,
		"core2_query_requestTimes_max_ms":     0,
		"core1_query_requestTimes_mean_ms":    0,
		"core1_update_totalTime_count":        3,
		"core1_query_serverErrors_count":      3,
		"core1_update_requestTimes_p99_ms":    297000000,
		"core2_query_totalTime_count":         3,
		"core2_update_requestTimes_max_ms":    0,
		"core2_query_requestTimes_p99_ms":     297000000,
		"core1_query_requestTimes_count":      3,
		"core1_query_requestTimes_median_ms":  0,
		"core1_query_clientErrors_count":      3,
		"core2_update_requestTimes_mean_ms":   0,
		"core2_update_requestTimes_median_ms": 0,
		"core2_update_requestTimes_p95_ms":    285000000,
		"core2_update_requestTimes_p999_ms":   2997000000,
		"core2_update_totalTime_count":        3,
		"core1_update_clientErrors_count":     3,
		"core2_query_serverErrors_count":      3,
		"core2_query_requests_count":          3,
		"core1_update_serverErrors_count":     3,
		"core1_update_requestTimes_p75_ms":    225000000,
		"core2_update_requestTimes_min_ms":    0,
		"core2_query_errors_count":            3,
		"core1_update_errors_count":           3,
		"core1_query_totalTime_count":         3,
		"core1_update_requestTimes_p95_ms":    285000000,
		"core2_query_requestTimes_p95_ms":     285000000,
		"core2_query_requestTimes_p999_ms":    2997000000,
		"core1_query_requestTimes_min_ms":     0,
		"core2_update_errors_count":           3,
		"core2_query_clientErrors_count":      3,
		"core1_update_requestTimes_min_ms":    0,
		"core1_query_requestTimes_max_ms":     0,
		"core1_query_requestTimes_p75_ms":     225000000,
		"core1_query_requestTimes_p999_ms":    2997000000,
		"core2_update_requestTimes_p75_ms":    225000000,
		"core2_update_timeouts_count":         3,
		"core1_query_requestTimes_p99_ms":     297000000,
		"core1_update_requests_count":         3,
		"core1_update_requestTimes_median_ms": 0,
		"core1_update_requestTimes_max_ms":    0,
		"core2_update_requestTimes_p99_ms":    297000000,
		"core2_query_requestTimes_count":      3,
		"core1_query_errors_count":            3,
		"core1_query_requests_count":          3,
		"core1_update_requestTimes_count":     3,
		"core2_update_requests_count":         3,
		"core2_query_requestTimes_median_ms":  0,
	}

	assert.Equal(t, expected, mod.Collect())
	assert.Equal(t, expected, mod.Collect())
}

func TestSolr_Collect_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.URL = ts.URL

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}
