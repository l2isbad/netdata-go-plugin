package nginxvts

import (
	"errors"
	"net/http"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func (vts *NginxVTS) validateConfig() error {
	if vts.URL == "" {
		return errors.New("URL not set")
	}

	if _, err := web.NewHTTPRequest(vts.Request); err != nil {
		return err
	}
	return nil
}

func (vts *NginxVTS) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(vts.Client)
}

func (vts *NginxVTS) initCharts() (*Charts, error) {
	charts := module.Charts{}

	if err := charts.Add(*nginxVtsMainCharts.Copy()...); err != nil {
		return nil, err
	}

	if err := charts.Add(*nginxVtsSharedZonesChart.Copy()...); err != nil {
		return nil, err
	}

	if err := charts.Add(*nginxVtsServerZonesCharts.Copy()...); err != nil {
		return nil, err
	}

	if len(charts) == 0 {
		return nil, errors.New("zero charts")
	}
	return &charts, nil
}
