package nvidia_nvml

import (
	"github.com/mindprince/gonvml"
	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		//DisabledByDefault: true,
		Create: func() modules.Module { return New() },
	}

	modules.Register("nvidia_nvml", creator)
}

// New creates NvidiaNVML with default values.
func New() *NvidiaNVML {
	return &NvidiaNVML{
		charts: &Charts{},
	}
}

// NvidiaNVML NvidiaNVML module.
type NvidiaNVML struct {
	modules.Base

	charts *Charts
}

// Cleanup makes cleanup.
func (NvidiaNVML) Cleanup() {
	_ = gonvml.Shutdown()
}

// Init makes initialization.
func (n NvidiaNVML) Init() bool {
	if err := gonvml.Initialize(); err != nil {
		n.Errorf("error on nvml initialization : %v", err)
		return false
	}

	gpus, err := getGPUs()

	if err != nil {
		n.Error(err)
		return false
	}

	for _, g := range gpus {
		_ = n.charts.Add(*createGPUCharts(g)...)
	}

	return true
}

// Check makes check.
func (NvidiaNVML) Check() bool {
	return true
}

// Charts creates Charts.
func (n NvidiaNVML) Charts() *Charts {
	return n.charts
}

// Collect collects metrics.
func (n *NvidiaNVML) Collect() map[string]int64 {
	m := make(map[string]int64)

	gpus, err := getGPUs()

	if err != nil {
		n.Error(err)
		return nil
	}

	for _, gpu := range gpus {
		for k, v := range gpu.stats.asMap() {
			m[gpu.uniqName()+"_"+k] = v
		}
	}

	return m
}
