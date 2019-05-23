package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/prometheus/prometheus/pkg/labels"
)

const (
	// defaults are cpu,cs,logical_disk,net,os,service,system,textfile
	collectorCPU      = "cpu"
	collectorLogDisks = "logical_disk"
	collectorNet      = "net"
	collectorOS       = "os"
	collectorSystem   = "system"
	collectorMemory   = "memory"

	metricCollectorDuration = "wmi_exporter_collector_duration_seconds"
	metricCollectorSuccess  = "wmi_exporter_collector_success"
)

func (w *WMI) collect() (map[string]int64, error) {
	scraped, err := w.prom.Scrape()
	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	w.collectScraped(mx, scraped)
	w.updateCharts(mx)

	return stm.ToMap(mx), nil
}

func (w *WMI) collectScraped(mx *metrics, scraped prometheus.Metrics) {
	collectCollection(mx, scraped)

	w.collectCPU(mx, scraped)
	w.collectOS(mx, scraped)
	w.collectMemory(mx, scraped)
	w.collectSystem(mx, scraped)
	w.collectNet(mx, scraped)
	w.collectLogicalDisk(mx, scraped)

	if mx.hasOS() && mx.hasMem() {
		v := sum(mx.OS.VisibleMemoryBytes, -mx.Memory.AvailableBytes)
		mx.Memory.UsedBytes = &v
	}
}

func collectCollection(mx *metrics, pms prometheus.Metrics) {
	mx.Collectors = &collectors{}
	collectCollectionDuration(mx, pms)
	collectCollectionSuccess(mx, pms)
}

func collectCollectionDuration(mx *metrics, pms prometheus.Metrics) {
	cr := newCollector("")
	for _, pm := range pms.FindByName(metricCollectorDuration) {
		name := pm.Labels.Get("collector")
		if name == "" {
			continue
		}
		if cr.ID != name {
			cr = mx.Collectors.get(name, true)
		}
		cr.Duration = pm.Value
	}
}

func collectCollectionSuccess(mx *metrics, pms prometheus.Metrics) {
	cr := newCollector("")
	for _, pm := range pms.FindByName(metricCollectorSuccess) {
		name := pm.Labels.Get("collector")
		if name == "" {
			continue
		}
		if cr.ID != name {
			cr = mx.Collectors.get(name, true)
		}
		cr.Success = pm.Value == 1
	}
}

func checkCollector(pms prometheus.Metrics, name string) (enabled, success bool) {
	m, err := labels.NewMatcher(labels.MatchEqual, "collector", name)
	if err != nil {
		panic(err)
	}

	pms = pms.FindByName(metricCollectorSuccess)
	ms := pms.Match(m)
	return ms.Len() > 0, ms.Max() == 1
}

func sum(vs ...float64) (s float64) {
	for _, v := range vs {
		s += v
	}
	return s
}
