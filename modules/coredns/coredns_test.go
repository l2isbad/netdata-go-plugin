package coredns

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testNoLoad, _   = ioutil.ReadFile("testdata/no_load.txt")
	testSomeLoad, _ = ioutil.ReadFile("testdata/some_load.txt")
)

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*CoreDNS)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestCoreDNS_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestCoreDNS_Cleanup(t *testing.T) { New().Cleanup() }

func TestCoreDNS_Init(t *testing.T) { assert.True(t, New().Init()) }

func TestCoreDNS_InitNG(t *testing.T) {
	job := New()
	job.URL = ""
	assert.False(t, job.Init())
}

func TestCoreDNS_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSomeLoad)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestCoreDNS_CheckNG(t *testing.T) {
	job := New()
	job.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestCoreDNS_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSomeLoad)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	require.True(t, job.Check())

	//for k, v := range job.Collect() {
	//	fmt.Println(fmt.Sprintf("\"%s\": %d,", k, v))
	//}

	expected := map[string]int64{
		"panic_count_total":                      0,
		"request_count_total":                    119,
		"request_count_by_type_total_A":          89,
		"request_count_by_type_total_AAAA":       29,
		"request_count_by_type_total_MX":         1,
		"response_count_by_rcode_total_NOERROR":  3,
		"response_count_by_rcode_total_SERVFAIL": 116,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestCoreDNS_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestCoreDNS_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}
