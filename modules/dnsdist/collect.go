package dnsdist

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

const (
	urlPathLocalStatistics = "/jsonstat"
)

func (d *DNSdist) collect() (map[string]int64, error) {
	statistics, err := d.scrapeStatistics()
	if err != nil {
		return nil, err;
	}

	collected := make(map[string]int64)

	d.collectStatistic(collected, statistics)
	return collected, nil
}

func (d *DNSdist) collectStatistic(collected map[string]int64, statistics *statisticMetrics) {
	for metric, value := range stm.ToMap(statistics) {
		collected[metric] = int64(value)
	}
}

func (d *DNSdist) scrapeStatistics() (*statisticMetrics, error) {
	req, err := web.NewHTTPRequest(d.Request)
	if err != nil {
		return nil, err
	}
	
	URL :=  []rune(d.Config.HTTP.Request.URL)
	req.URL.Opaque = string(URL[5:]) + urlPathLocalStatistics

	var statistics statisticMetrics
	if err := d.doOKDecode(req, &statistics); err != nil {
		return nil, err
	}

	return &statistics, nil
}

func (d *DNSdist) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := d.httpClient.Do(req)
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

