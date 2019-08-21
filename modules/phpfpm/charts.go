package phpfpm

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "connections",
		Title: "PHP-FPM Active Connections",
		Units: "connections",
		Fam:   "active connections",
		Ctx:   "phpfpm.connections",
		Dims: Dims{
			{ID: "active"},
			{ID: "maxActive", Name: "max active"},
			{ID: "idle"},
		},
	},
	{
		ID:    "requests",
		Title: "PHP-FPM Requests",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "phpfpm.requests",
		Dims: Dims{
			{ID: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "performance",
		Title: "PHP-FPM Performance",
		Units: "status",
		Fam:   "performance",
		Ctx:   "phpfpm.performance",
		Dims: Dims{
			{ID: "reached", Name: "max children reached"},
			{ID: "slow", Name: "slow requests"},
		},
	},
	{
		ID:    "request_duration",
		Title: "PHP-FPM Request Duration",
		Units: "milliseconds",
		Fam:   "request duration",
		Ctx:   "phpfpm.request_duration",
		Dims: Dims{
			{ID: "minReqDur", Name: "min", Div: 1000},
			{ID: "maxReqDur", Name: "max", Div: 1000},
			{ID: "avgReqDur", Name: "avg", Div: 1000},
		},
	},
	{
		ID:    "request_cpu",
		Title: "PHP-FPM Request CPU",
		Units: "percentage",
		Fam:   "request CPU",
		Ctx:   "phpfpm.request_cp",
		Dims: Dims{
			{ID: "minReqCpu", Name: "min"},
			{ID: "maxReqCpu", Name: "max"},
			{ID: "avgReqCpu", Name: "avg"},
		},
	},
	{
		ID:    "request_mem",
		Title: "PHP-FPM Request Memory",
		Units: "KB",
		Fam:   "request memory",
		Ctx:   "phpfpm.request_mem",
		Dims: Dims{
			{ID: "minReqMem", Name: "min", Div: 1024},
			{ID: "maxReqMem", Name: "max", Div: 1024},
			{ID: "avgReqMem", Name: "avg", Div: 1024},
		},
	},
}
