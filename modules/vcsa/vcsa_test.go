package vcsa

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testNewVCSA() *VCSA {
	vc := New()
	vc.UserURL = "https://127.0.0.1:38001"
	vc.Username = "user"
	vc.Password = "pass"
	return vc
}

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*VCSA)(nil), job)
}

func TestVCSA_Init(t *testing.T) {
	job := testNewVCSA()

	assert.True(t, job.Init())
	assert.NotNil(t, job.client)
}

func TestVCenter_InitErrorOnValidatingInitParameters(t *testing.T) {
	job := New()

	assert.False(t, job.Init())
}

func TestVCenter_InitErrorOnCreatingClient(t *testing.T) {
	job := testNewVCSA()
	job.ClientTLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestVCenter_Check(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	job.client = &mockVCenterHealthClient{}

	assert.True(t, job.Check())
}

func TestVCenter_CheckErrorOnLogin(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	job.client = &mockVCenterHealthClient{
		login: func() error { return errors.New("login mock error") },
	}

	assert.False(t, job.Check())
}

func TestVCenter_CheckEnsureLoggedIn(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	mock := &mockVCenterHealthClient{}
	job.client = mock

	assert.True(t, job.Check())
	assert.True(t, mock.loginCalls == 1)
}

func TestVCenter_Cleanup(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	mock := &mockVCenterHealthClient{}
	job.client = mock
	job.Cleanup()

	assert.True(t, mock.logoutCalls == 1)
}

func TestVCenter_CleanupWithNilClient(t *testing.T) {
	job := testNewVCSA()

	assert.NotPanics(t, job.Cleanup)
}

func TestVCenter_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestVCenter_Collect(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	mock := &mockVCenterHealthClient{}
	job.client = mock

	expected := map[string]int64{
		"appl_mgmt":         0,
		"database_storage":  0,
		"load":              0,
		"mem":               0,
		"software_packages": 0,
		"storage":           0,
		"swap":              0,
		"system":            0,
	}
	assert.Equal(t, expected, job.Collect())
}

func TestVCenter_CollectEnsurePingIsCalled(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	mock := &mockVCenterHealthClient{}
	job.client = mock
	job.Collect()

	assert.True(t, mock.pingCalls == 1)
}

func TestVCenter_CollectErrorOnPing(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	mock := &mockVCenterHealthClient{
		ping: func() error { return errors.New("ping mock error") },
	}
	job.client = mock

	assert.Zero(t, job.Collect())
}

func TestVCenter_CollectErrorOnHealthCalls(t *testing.T) {
	job := testNewVCSA()
	require.True(t, job.Init())
	mock := &mockVCenterHealthClient{
		applMgmt:         func() (string, error) { return "", errors.New("applMgmt mock error") },
		databaseStorage:  func() (string, error) { return "", errors.New("databaseStorage mock error") },
		load:             func() (string, error) { return "", errors.New("load mock error") },
		mem:              func() (string, error) { return "", errors.New("mem mock error") },
		softwarePackages: func() (string, error) { return "", errors.New("softwarePackages mock error") },
		storage:          func() (string, error) { return "", errors.New("storage mock error") },
		swap:             func() (string, error) { return "", errors.New("swap mock error") },
		system:           func() (string, error) { return "", errors.New("system mock error") },
	}
	job.client = mock

	assert.Zero(t, job.Collect())
}

type mockVCenterHealthClient struct {
	login            func() error
	logout           func() error
	ping             func() error
	applMgmt         func() (string, error)
	databaseStorage  func() (string, error)
	load             func() (string, error)
	mem              func() (string, error)
	softwarePackages func() (string, error)
	storage          func() (string, error)
	swap             func() (string, error)
	system           func() (string, error)
	loginCalls       int
	logoutCalls      int
	pingCalls        int
}

func (m *mockVCenterHealthClient) Login() error {
	m.loginCalls += 1
	if m.login == nil {
		return nil
	}
	return m.login()
}

func (m *mockVCenterHealthClient) Logout() error {
	m.logoutCalls += 1
	if m.logout == nil {
		return nil
	}
	return m.logout()
}

func (m *mockVCenterHealthClient) Ping() error {
	m.pingCalls += 1
	if m.ping == nil {
		return nil
	}
	return m.ping()
}

func (m mockVCenterHealthClient) ApplMgmt() (string, error) {
	if m.applMgmt == nil {
		return "green", nil
	}
	return m.applMgmt()
}

func (m mockVCenterHealthClient) DatabaseStorage() (string, error) {
	if m.databaseStorage == nil {
		return "green", nil
	}
	return m.databaseStorage()
}

func (m mockVCenterHealthClient) Load() (string, error) {
	if m.load == nil {
		return "green", nil
	}
	return m.load()
}

func (m mockVCenterHealthClient) Mem() (string, error) {
	if m.mem == nil {
		return "green", nil
	}
	return m.mem()
}

func (m mockVCenterHealthClient) SoftwarePackages() (string, error) {
	if m.softwarePackages == nil {
		return "green", nil
	}
	return m.softwarePackages()
}

func (m mockVCenterHealthClient) Storage() (string, error) {
	if m.storage == nil {
		return "green", nil
	}
	return m.storage()
}

func (m mockVCenterHealthClient) Swap() (string, error) {
	if m.swap == nil {
		return "green", nil
	}
	return m.swap()
}

func (m mockVCenterHealthClient) System() (string, error) {
	if m.system == nil {
		return "green", nil
	}
	return m.system()
}