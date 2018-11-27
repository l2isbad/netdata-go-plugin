package modules

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockModule struct {
	Base
	initFunc      func() bool
	checkFunc     func() bool
	getChartsFunc func() *Charts
	getDataDunc   func() map[string]int64
}

func (m mockModule) Init() bool {
	return m.initFunc()
}

func (m mockModule) Check() bool {
	return m.checkFunc()
}

func (m mockModule) GetCharts() *Charts {
	return m.getChartsFunc()
}

func (m mockModule) GetData() map[string]int64 {
	return m.getDataDunc()
}

func TestNewJob(t *testing.T) {
	assert.IsType(
		t,
		(*Job)(nil),
		NewJob("example", nil, ioutil.Discard, nil),
	)
}

func TestJob_FullName(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.FullName(), "modName_modName")
	job.Nam = "jobName"
	assert.Equal(t, job.FullName(), "modName_jobName")

}

func TestJob_ModuleName(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.ModuleName(), "modName")
}

func TestJob_Name(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.Name(), "modName")
	job.Nam = "jobName"
	assert.Equal(t, job.Name(), "jobName")
}

func TestJob_Initialized(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.Initialized(), job.initialized)
	job.initialized = true
	assert.Equal(t, job.Initialized(), job.initialized)

}

func TestJob_Panicked(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.Panicked(), job.panicked)
	job.panicked = true
	assert.Equal(t, job.Panicked(), job.panicked)

}

func TestJob_AutoDetectionRetry(t *testing.T) {
	job := NewJob("modName", &mockModule{}, ioutil.Discard, nil)

	assert.Equal(t, job.AutoDetectionRetry(), job.AutoDetectRetry)
	job.AutoDetectRetry = 1
	assert.Equal(t, job.AutoDetectionRetry(), job.AutoDetectRetry)

}

func TestJob_Init(t *testing.T) {
	okJob := NewJob(
		"modName",
		&mockModule{initFunc: func() bool { return true }},
		ioutil.Discard, nil,
	)

	assert.True(t, okJob.Init())
	assert.True(t, okJob.Initialized())

	panicJob := NewJob(
		"modName",
		&mockModule{initFunc: func() bool { panic("panic in init") }},
		ioutil.Discard, nil,
	)

	assert.False(t, panicJob.Init())
	assert.False(t, panicJob.Initialized())
}

func TestJob_Check(t *testing.T) {
	okJob := NewJob(
		"modName",
		&mockModule{checkFunc: func() bool { return true }},
		ioutil.Discard, nil,
	)

	assert.True(t, okJob.Check())

	panicJob := NewJob(
		"modName",
		&mockModule{checkFunc: func() bool { panic("panic in test") }},
		ioutil.Discard, nil,
	)

	assert.False(t, panicJob.Check())
}

func TestJob_PostCheck(t *testing.T) {

}

func TestJob_Start(t *testing.T) {

}

func TestJob_Stop(t *testing.T) {

}

func TestJob_Tick(t *testing.T) {

}

func TestJob_MainLoop(t *testing.T) {

}
