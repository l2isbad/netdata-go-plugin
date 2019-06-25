package phpdaemon

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "total_workers",
		Title: "Total Workers",
		Units: "workers",
		Fam:   "workers",
		Ctx:   "phpdaemon.total_workers",
		Dims: Dims{
			{ID: "total"},
		},
	},
	{
		ID:    "workers",
		Title: "Workers State",
		Units: "workers",
		Fam:   "workers",
		Ctx:   "phpdaemon.workers",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "alive"},
			{ID: "shutdown"},
		},
	},
	{
		ID:    "alive_workers",
		Title: "Alive Workers State",
		Units: "workers",
		Fam:   "workers",
		Ctx:   "phpdaemon.alive_workers",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "idle"},
			{ID: "busy"},
			{ID: "reloading"},
		},
	},
	{
		ID:    "idle_workers",
		Title: "Idle Workers State",
		Units: "workers",
		Fam:   "workers",
		Ctx:   "phpdaemon.idle_workers",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "preinit"},
			{ID: "init"},
			{ID: "initialized"},
		},
	},
}

var uptimeChart = Chart{
	ID:    "uptime",
	Title: "Uptime",
	Units: "seconds",
	Fam:   "uptime",
	Ctx:   "phpdaemon.uptime",
	Dims: Dims{
		{ID: "uptime", Name: "time"},
	},
}
