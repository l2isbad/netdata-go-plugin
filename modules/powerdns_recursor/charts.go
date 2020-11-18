package powerdns_recursor

import "github.com/netdata/go.d.plugin/agent/module"

// https://doc.powerdns.com/recursor/metrics.html
// https://docs.powerdns.com/recursor/performance.html#recursor-caches

var charts = module.Charts{
	{
		ID:    "questions_in",
		Title: "Questions In",
		Units: "questions",
		Fam:   "questions",
		Ctx:   "powerdns_recursor.questions_in",
		Dims: module.Dims{
			{ID: "questions", Name: "total", Algo: module.Incremental},
			{ID: "tcp-questions", Name: "tcp", Algo: module.Incremental},
			{ID: "ipv6-questions", Name: "ipv6", Algo: module.Incremental},
		},
	},
	{
		ID:    "questions_out",
		Title: "Questions Out",
		Units: "questions",
		Fam:   "questions",
		Ctx:   "powerdns_recursor.questions_out",
		Dims: module.Dims{
			{ID: "all-outqueries", Name: "udp", Algo: module.Incremental},
			{ID: "tcp-outqueries", Name: "tcp", Algo: module.Incremental},
			{ID: "ipv6-outqueries", Name: "ipv6", Algo: module.Incremental},
			{ID: "throttled-outqueries", Name: "throttled", Algo: module.Incremental},
		},
	},
	{
		ID:    "answer_time",
		Title: "Queries answered within a time range",
		Units: "queries",
		Fam:   "performance",
		Ctx:   "powerdns_recursor.answer_time",
		Dims: module.Dims{
			{ID: "answers0-1", Name: "0-1ms", Algo: module.Incremental},
			{ID: "answers1-10", Name: "1-10ms", Algo: module.Incremental},
			{ID: "answers10-100", Name: "10-100ms", Algo: module.Incremental},
			{ID: "answers100-1000", Name: "100-1000ms", Algo: module.Incremental},
			{ID: "answers-slow", Name: "slow", Algo: module.Incremental},
		},
	},
	{
		ID:    "timeouts",
		Title: "Timeouts on outgoing UDP queries since starting",
		Units: "timeouts",
		Fam:   "performance",
		Ctx:   "powerdns_recursor.timeouts",
		Dims: module.Dims{
			{ID: "outgoing-timeouts", Name: "total", Algo: module.Incremental},
			{ID: "outgoing4-timeouts", Name: "ipv4", Algo: module.Incremental},
			{ID: "outgoing6-timeouts", Name: "ipv6", Algo: module.Incremental},
		},
	},
	{
		ID:    "question_drops",
		Title: "Dropped questions",
		Units: "questions",
		Fam:   "performance",
		Ctx:   "powerdns_recursor.question_drops",
		Dims: module.Dims{
			{ID: "over-capacity-drops", Algo: module.Incremental},
			{ID: "query-pipe-full-drops", Algo: module.Incremental},
			{ID: "too-old-drops", Algo: module.Incremental},
			{ID: "truncated-drops", Algo: module.Incremental},
			{ID: "empty-queries", Algo: module.Incremental},
		},
	},
	{
		ID:    "cache_usage",
		Title: "Cache Usage",
		Units: "events/s",
		Fam:   "cache",
		Ctx:   "powerdns_recursor.cache_usage",
		Dims: module.Dims{
			{ID: "cache-hits", Algo: module.Incremental},
			{ID: "cache-misses", Algo: module.Incremental},
			{ID: "packetcache-hits", Name: "packet-cache-hits", Algo: module.Incremental},
			{ID: "packetcache-misses", Name: "packet-cache-misses", Algo: module.Incremental},
		},
	},
	{
		ID:    "cache_size",
		Title: "Cache Size",
		Units: "entries",
		Fam:   "cache",
		Ctx:   "powerdns_recursor.cache_size",
		Dims: module.Dims{
			{ID: "cache-entries", Name: "cache"},
			{ID: "packetcache-entries", Name: "packet-cache"},
			{ID: "negcache-entries", Name: "negative-cache"},
		},
	},
}
