package modules

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/l2isbad/go.d.plugin/logger"
	"github.com/l2isbad/go.d.plugin/modules/internal/apiwriter"
	"github.com/l2isbad/go.d.plugin/modules/internal/chartstate"
)

type switchResult string

const (
	switchContinue switchResult = "continue"
	switchSkip     switchResult = "skip"
	switchBreak    switchResult = "break"
)

type JobConfig struct {
	JobName            string `yaml:"job_name"`
	OverrideName       string `yaml:"name"`
	UpdateEvery        int    `yaml:"update_every" validate:"gte=1"`
	AutoDetectionRetry int    `yaml:"autodetection_retry" validate:"gte=0"`
	ChartCleanup       int    `yaml:"chart_cleanup" validate:"gte=0"`
	MaxRetries         int    `yaml:"retries" validate:"gte=0"`
}

// JobNewConfig returns JobConfig with default values
func JobNewConfig() *JobConfig {
	return &JobConfig{
		UpdateEvery:        1,
		AutoDetectionRetry: 0,
		ChartCleanup:       10,
		MaxRetries:         60,
	}
}

type Job struct {
	*JobConfig
	*logger.Logger
	module Module

	Inited     bool
	Panicked   bool
	ModuleName string

	charts       *Charts
	tick         chan int
	shutdownHook chan struct{}
	out          io.Writer
	buf          *bytes.Buffer
	apiWriter    apiwriter.APIWriter

	priority int
	retries  int
	prevRun  time.Time
}

func (j Job) Name() string {
	if j.ModuleName == j.JobName {
		return j.ModuleName
	}
	return fmt.Sprintf("%s_%s", j.ModuleName, j.JobName)
}

func NewJob(modName string, module Module, config *JobConfig, out io.Writer) *Job {
	buf := &bytes.Buffer{}
	return &Job{
		ModuleName:   modName,
		module:       module,
		JobConfig:    config,
		out:          out,
		tick:         make(chan int),
		shutdownHook: make(chan struct{}),
		buf:          buf,
		priority:     70000,
		apiWriter:    apiwriter.APIWriter{Writer: buf},
	}
}

func (j *Job) Init() error {
	j.Logger = logger.New(j.ModuleName, j.JobName)
	j.module.SetUpdateEvery(j.UpdateEvery)
	j.module.SetModuleName(j.ModuleName)
	j.module.SetLogger(j.Logger)

	return j.module.Init()
}

func (j *Job) Check() bool {
	defer func() {
		if r := recover(); r != nil {
			j.Panicked = true
			j.Errorf("PANIC %v", r)
		}

	}()
	return j.module.Check()
}

func (j *Job) PostCheck() bool {
	j.UpdateEvery = j.module.UpdateEvery()
	j.ModuleName = j.module.ModuleName()
	logger.SetModName(j.Logger, j.ModuleName)

	charts := j.module.GetCharts()
	if charts == nil {
		j.Error("GetCharts() [FAILED]")
		return false
	}

	j.charts = charts
	return true
}

func (j *Job) Tick(clock int) {
	select {
	case j.tick <- clock:
	default:
		j.Errorf("Skip the tick due to previous run hasn't been finished.")
	}
}

func (j *Job) MainLoop() {
LOOP:
	for {
		select {
		case <-j.shutdownHook:
			break LOOP
		case t := <-j.tick:
			if t%j.UpdateEvery != 0 {
				continue LOOP
			}
			j.Info(11111111111)
		}

		//curTime := time.Now()
		//if j.prevRun.IsZero() {
		//	sinceLast := 0
		//} else {
		//	sinceLast := convertTo(curTime.Sub(j.prevRun), time.Microsecond)
		//}
		//
		//data := j.getData()
		//
		//if data == nil {
		//	j.retries++
		//	continue
		//}
		//j.buf.Reset()
		//// TODO write data
		//io.Copy(j.out, j.buf)
	}
}

func (j *Job) Shutdown() {
	select {
	case j.shutdownHook <- struct{}{}:
	default:
	}
}

func (j *Job) getData() (result map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC: %v", r)
			j.Panicked = true
		}
	}()
	return j.module.GetData()
}

func (j *Job) AutoDetectionRetry() int {
	return j.JobConfig.AutoDetectionRetry
}

func (j *Job) PopulateMetrics(data map[string]int64, sinceLast int) bool {
	var updated int
LOOP:
	for _, chart := range *j.charts {
		for {
			switch j.switchChartState(chart, data) {
			case switchContinue:
				continue
			case switchSkip:
				continue LOOP
			case switchBreak:
				break
			}
		}
		if data == nil {
			continue
		}
		j.apiWriter.Begin("typeName", chart.ID, sinceLast)
		var chartUpdated bool

		for _, dim := range chart.Dims {
			if v, ok := data[dim.ID]; ok {
				j.apiWriter.Set(dim.ID, v)
				chartUpdated = true
			}
		}

		for _, variable := range chart.Vars {
			if v, ok := data[variable.ID]; ok {
				j.apiWriter.Set(variable.ID, v)
			}
		}

		if chartUpdated {
			updated++
			chart.retries = 0
		} else if chart.retries++; j.ChartCleanup > 0 && chart.retries >= j.ChartCleanup {
			chart.state = chartstate.MarkedObsolete
		}

		j.apiWriter.End()
	}

	return updated > 0
}

func (j *Job) switchChartState(chart *Chart, data map[string]int64) switchResult {
	switch chart.state {
	case chartstate.Initial:
		chart.priority = j.priority
		j.priority++
		chart.state = chartstate.New
	case chartstate.New:
		chart.state = chartstate.Created
	case chartstate.Created:
		return switchBreak
	case chartstate.Rebuilt:

		chart.state = chartstate.New
	case chartstate.Recovered:

		chart.state = chartstate.New
	case chartstate.MarkedObsolete:

		chart.state = chartstate.Obsoleted
	case chartstate.Obsoleted:
		if canChartBeUpdated(chart, data) {
			chart.state = chartstate.New
		} else {
			return switchSkip
		}
		chart.state = chartstate.New
	case chartstate.MarkedRemove:
		chart.state = chartstate.MarkedDelete
	case chartstate.MarkedDelete:
		return switchSkip
	}
	return switchContinue
}

func convertTo(from time.Duration, to time.Duration) int {
	return int(int64(from) / (int64(to) / int64(time.Nanosecond)))
}

func canChartBeUpdated(chart *Chart, data map[string]int64) bool {
	for _, dim := range chart.Dims {
		if _, ok := data[dim.ID]; ok {
			return true
		}
	}
	return false
}
