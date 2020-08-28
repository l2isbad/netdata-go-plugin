package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathNodesLocalStats = "/_nodes/_local/stats"
	urlPathClusterHealth   = "/_cluster/health"
	urlPathClusterStats    = "/_cluster/stats"
)

func (es *Elasticsearch) collect() (map[string]int64, error) {
	mx := es.scrapeElasticsearch()

	return stm.ToMap(mx), nil
}

func (es *Elasticsearch) scrapeElasticsearch() *esMetrics {
	var mx esMetrics
	wg := &sync.WaitGroup{}
	for _, task := range []func(mx *esMetrics){
		es.scrapeLocalNodeStats,
		es.scrapeClusterHealth,
		es.scrapeClusterStats,
	} {
		wg.Add(1)
		task := task
		go func() { defer wg.Done(); task(&mx) }()
	}
	wg.Wait()
	return &mx
}

func (es *Elasticsearch) scrapeLocalNodeStats(mx *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathNodesLocalStats

	var stats struct{ Nodes []esNodeStats }
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	if len(stats.Nodes) > 0 {
		mx.LocalNodeStats = &stats.Nodes[0]
	}
}

func (es *Elasticsearch) scrapeClusterHealth(mx *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathClusterHealth

	var health esClusterHealth
	if err := es.doOKDecode(req, &health); err != nil {
		es.Warning(err)
		return
	}
	mx.ClusterHealth = &health
}

func (es *Elasticsearch) scrapeClusterStats(mx *esMetrics) {
	req, _ := web.NewHTTPRequest(es.Request)
	req.URL.Path = urlPathClusterStats

	var stats esClusterStats
	if err := es.doOKDecode(req, &stats); err != nil {
		es.Warning(err)
		return
	}
	mx.ClusterStats = &stats
}

func (es *Elasticsearch) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := es.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on HTTP request '%s': %v", req.URL, err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("'%s' returned HTTP status code: %d", req.URL, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(in); err != nil {
		return fmt.Errorf("error on decoding response from '%s': %v", req.URL, err)
	}

	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
