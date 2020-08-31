package elasticsearch

import "github.com/netdata/go-orchestrator/module"

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
)

// TODO: indices operations charts: query_latency, fetch_latency, index_latency, flush_latency

var (
	nodeStatsIndicesIndexingCharts = Charts{
		{
			ID:    "indices_index_total",
			Title: "Index Operations Total",
			Units: "operations/s",
			Fam:   "indices index",
			Ctx:   "elasticsearch.indices_index_total",
			Dims: Dims{
				{ID: "indices_indexing_index_total", Name: "index", Algo: module.Incremental},
			},
		},
		{
			ID:    "indices_index_current",
			Title: "Index Operations Current",
			Units: "operations",
			Fam:   "indices index",
			Ctx:   "elasticsearch.indices_index_current",
			Dims: Dims{
				{ID: "indices_indexing_index_current", Name: "index"},
			},
		},
		{
			ID:    "indices_index_time",
			Title: "Time Spent On Index Operations",
			Units: "milliseconds",
			Fam:   "indices index",
			Ctx:   "elasticsearch.indices_index_time",
			Dims: Dims{
				{ID: "indices_indexing_index_time_in_millis", Name: "query", Algo: module.Incremental},
			},
		},
	}
	nodeStatsIndicesSearchCharts = Charts{
		{
			ID:    "indices_search_total",
			Title: "Search Operations Total",
			Units: "operations/s",
			Fam:   "indices search",
			Ctx:   "elasticsearch.indices_search_total",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "indices_search_query_total", Name: "queries", Algo: module.Incremental},
				{ID: "indices_search_fetch_total", Name: "fetches", Algo: module.Incremental},
			},
		},
		{
			ID:    "indices_search_current",
			Title: "Search Operations Current",
			Units: "events/s",
			Fam:   "indices search",
			Ctx:   "elasticsearch.indices_search_current",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "indices_search_query_current", Name: "queries"},
				{ID: "indices_search_fetch_current", Name: "fetches"},
			},
		},
		{
			ID:    "indices_search_time",
			Title: "Time Spent On Search Operations",
			Units: "milliseconds",
			Fam:   "indices search",
			Ctx:   "elasticsearch.indices_search_time",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "indices_search_query_time_in_millis", Name: "query", Algo: module.Incremental},
				{ID: "indices_search_fetch_time_in_millis", Name: "fetch", Algo: module.Incremental},
			},
		},
		// TODO: search_latency
	}
	nodeStatsIndicesRefreshCharts = Charts{
		{
			ID:    "indices_refresh_total",
			Title: "Refresh Operations Total",
			Units: "operations/s",
			Fam:   "indices refresh",
			Ctx:   "elasticsearch.indices_refresh_total",
			Dims: Dims{
				{ID: "indices_refresh_total", Name: "refresh", Algo: module.Incremental},
			},
		},
		{
			ID:    "indices_refresh_time",
			Title: "Time Spent On Refresh Operations",
			Units: "milliseconds",
			Fam:   "indices refresh",
			Ctx:   "elasticsearch.indices_refresh_time",
			Dims: Dims{
				{ID: "indices_refresh_total_time_in_millis", Name: "refresh", Algo: module.Incremental},
			},
		},
	}
	nodeStatsIndicesFlushCharts = Charts{
		{
			ID:    "indices_flush_total",
			Title: "Flush Operations Total",
			Units: "operations/s",
			Fam:   "indices flush",
			Ctx:   "elasticsearch.indices_flush_total",
			Dims: Dims{
				{ID: "indices_flush_total", Name: "flush", Algo: module.Incremental},
			},
		},
		{
			ID:    "indices_flush_time",
			Title: "Time Spent On Flush Operations",
			Units: "milliseconds",
			Fam:   "indices flush",
			Ctx:   "elasticsearch.indices_flush_time",
			Dims: Dims{
				{ID: "indices_flush_total_time_in_millis", Name: "flush", Algo: module.Incremental},
			},
		},
	}
	nodeStatsIndicesFielddataCharts = Charts{
		{
			ID:    "indices_fielddata_memory_usage",
			Title: "Fielddata Cache Memory Usage",
			Units: "bytes",
			Fam:   "indices fielddata",
			Ctx:   "elasticsearch.indices_fielddata_memory_usage",
			Type:  module.Area,
			Dims: Dims{
				{ID: "indices_fielddata_memory_size_in_bytes", Name: "total"},
			},
		},
		{
			ID:    "indices_fielddata_evictions",
			Title: "Fielddata Evictions",
			Units: "evictions/s",
			Fam:   "indices fielddata",
			Ctx:   "elasticsearch.indices_fielddata_evictions",
			Dims: Dims{
				{ID: "indices_fielddata_evictions", Name: "evictions", Algo: module.Incremental},
			},
		},
	}
	nodeStatsIndicesSegmentsCharts = Charts{
		{
			ID:    "indices_segments_count",
			Title: "Segments Count",
			Units: "num",
			Fam:   "indices segments",
			Ctx:   "elasticsearch.indices_segments_count",
			Dims: Dims{
				{ID: "indices_segments_count", Name: "segments"},
			},
		},
		{
			ID:    "indices_segments_memory_usage_total",
			Title: "Segments Memory Usage Total",
			Units: "bytes",
			Fam:   "indices segments",
			Ctx:   "elasticsearch.indices_segments_memory_usage_total",
			Dims: Dims{
				{ID: "indices_memory_in_bytes", Name: "total"},
			},
		},
		{
			ID:    "indices_segments_memory_usage",
			Title: "Indices Segments Memory Usage",
			Units: "bytes",
			Fam:   "indices segments",
			Ctx:   "elasticsearch.indices_segments_memory_usage",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "indices_segments_terms_memory_in_bytes", Name: "terms"},
				{ID: "indices_segments_stored_fields_memory_in_bytes", Name: "stored_fields"},
				{ID: "indices_segments_term_vectors_memory_in_bytes", Name: "term_vectors"},
				{ID: "indices_segments_norms_memory_in_byte", Name: "norms"},
				{ID: "indices_segments_points_memory_in_bytes", Name: "points"},
				{ID: "indices_segments_doc_values_memory_in_bytes", Name: "doc_values"},
				{ID: "indices_segments_index_writer_memory_in_bytes", Name: "index_writer"},
				{ID: "indices_segments_version_map_memory_in_bytes", Name: "version_map"},
				{ID: "indices_segments_fixed_bit_set_memory_in_bytes", Name: "fixed_bit_set"},
			},
		},
	}
	nodeStatsIndicesTranslogCharts = Charts{
		{
			ID:    "indices_translog_operations",
			Title: "Translog Operations",
			Units: "operations",
			Fam:   "indices translog",
			Ctx:   "elasticsearch.indices_translog_operations",
			Type:  module.Area,
			Dims: Dims{
				{ID: "indices_translog_operations", Name: "total"},
				{ID: "indices_translog_uncommitted_operations", Name: "uncommitted"},
			},
		},
		{
			ID:    "index_translog_size",
			Title: "Translog Size",
			Units: "bytes",
			Fam:   "indices translog",
			Ctx:   "elasticsearch.indices_translog_size",
			Type:  module.Area,
			Dims: Dims{
				{ID: "indices_translog_size_in_bytes", Name: "total"},
				{ID: "indices_translog_uncommitted_size_in_bytes", Name: "uncommitted"},
			},
		},
	}
	nodeStatsProcessCharts = Charts{
		{
			ID:    "file_descriptors",
			Title: "File Descriptors",
			Units: "num",
			Fam:   "process",
			Ctx:   "elasticsearch.file_descriptors",
			Dims: Dims{
				{ID: "node_stats_process_open_file_descriptors", Name: "open"},
			},
			Vars: Vars{
				{ID: "node_stats_process_max_file_descriptors"},
			},
		},
	}
	nodeStatsJVMCharts = Charts{
		{
			ID:    "jvm_mem_heap",
			Title: "JVM Heap Percentage Currently in Use",
			Units: "percentage",
			Fam:   "jvm",
			Ctx:   "elasticsearch.jvm_heap",
			Type:  module.Area,
			Dims: Dims{
				{ID: "jvm_mem_heap_used_percent", Name: "inuse"},
			},
		},
		{
			ID:    "jvm_mem_heap_bytes",
			Title: "JVM Heap Commit And Usage",
			Units: "bytes",
			Fam:   "jvm",
			Ctx:   "elasticsearch.jvm_heap_bytes",
			Type:  module.Area,
			Dims: Dims{
				{ID: "jvm_mem_heap_committed_in_bytes", Name: "committed"},
				{ID: "jvm_mem_heap_used_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "jvm_buffer_pool_count",
			Title: "JVM Buffers",
			Units: "pools",
			Fam:   "jvm",
			Ctx:   "elasticsearch.jvm_buffer_pool_count",
			Dims: Dims{
				{ID: "jvm_buffer_pools_direct_count", Name: "direct"},
				{ID: "jvm_buffer_pools_mapped_count", Name: "mapped"},
			},
		},
		{
			ID:    "jvm_direct_buffers_memory",
			Title: "JVM Direct Buffers Memory",
			Units: "bytes",
			Fam:   "jvm",
			Ctx:   "elasticsearch.jvm_direct_buffers_memory",
			Type:  module.Area,
			Dims: Dims{
				{ID: "jvm_buffer_pools_direct_total_capacity_in_bytes", Name: "total"},
				{ID: "jvm_buffer_pools_direct_used_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "jvm_mapped_buffers_memory",
			Title: "JVM Mapped Buffers Memory",
			Units: "bytes",
			Fam:   "jvm",
			Ctx:   "elasticsearch.jvm_mapped_buffers_memory",
			Type:  module.Area,
			Dims: Dims{
				{ID: "jvm_buffer_pools_mapped_total_capacity_in_bytes", Name: "total"},
				{ID: "jvm_buffer_pools_mapped_used_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "jvm_gc_count",
			Title: "Garbage Collections",
			Units: "events/s",
			Fam:   "jvm",
			Ctx:   "elasticsearch.gc_count",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "jvm_gc_collectors_young_collection_count", Name: "young", Algo: module.Incremental},
				{ID: "jvm_gc_collectors_old_collection_count", Name: "old", Algo: module.Incremental},
			},
		},
		{
			ID:    "jvm_gc_time",
			Title: "Time Spent On Garbage Collections",
			Units: "milliseconds",
			Fam:   "jvm",
			Ctx:   "elasticsearch.gc_time",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "jvm_gc_collectors_young_collection_time_in_millis", Name: "young", Algo: module.Incremental},
				{ID: "jvm_gc_collectors_old_collection_time_in_millis", Name: "old", Algo: module.Incremental},
			},
		},
	}
	nodeStatsThreadPoolCharts = Charts{
		{
			ID:    "thread_pool_queued",
			Title: "Queued Threads In Thread Pool",
			Units: "num",
			Fam:   "thread pool",
			Ctx:   "elasticsearch.thread_pool_queued",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "thread_pool_write_queue", Name: "write"},
				{ID: "thread_pool_search_queue", Name: "search"},
			},
		},
		{
			ID:    "thread_pool_rejected",
			Title: "Rejected Threads In Thread Pool",
			Units: "num",
			Fam:   "thread pool",
			Ctx:   "elasticsearch.thread_pool_rejected",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "thread_pool_write_rejected", Name: "write"},
				{ID: "thread_pool_search_rejected", Name: "search"},
			},
		},
	}
	nodeStatsTransportCharts = Charts{
		{
			ID:    "cluster_communication_packets",
			Title: "Cluster Communication",
			Units: "pps",
			Fam:   "transport",
			Ctx:   "elasticsearch.cluster_communication_packets",
			Dims: Dims{
				{ID: "node_stats_transport_rx_count", Name: "rx", Algo: module.Incremental},
				{ID: "node_stats_transport_tx_count", Name: "tx", Mul: -1, Algo: module.Incremental},
			},
		},
		{
			ID:    "cluster_communication",
			Title: "Cluster Communication",
			Units: "bytes",
			Fam:   "transport",
			Ctx:   "elasticsearch.cluster_communication",
			Dims: Dims{
				{ID: "node_stats_transport_rx_size_in_bytes", Name: "rx", Algo: module.Incremental},
				{ID: "node_stats_transport_tx_size_in_bytes", Name: "tx", Mul: -1, Algo: module.Incremental},
			},
		},
	}
	nodeStatsHTTPCharts = Charts{
		{
			ID:    "http_connections",
			Title: "HTTP Connections",
			Units: "connections",
			Fam:   "http",
			Ctx:   "elasticsearch.http_connections",
			Dims: Dims{
				{ID: "node_stats_http_current_open", Name: "open"},
			},
		},
	}
	nodeStatsBreakersCharts = Charts{
		{
			ID:    "breakers_trips",
			Title: "Circuit Breaker Trips Count",
			Units: "trips/s",
			Fam:   "circuit breakers",
			Ctx:   "elasticsearch.breakers_trips",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_breakers_requests_tripped", Name: "requests", Algo: module.Incremental},
				{ID: "node_breakers_fielddata_tripped", Name: "fielddata", Algo: module.Incremental},
				{ID: "node_breakers_in_flight_requests_tripped", Name: "in_flight_requests", Algo: module.Incremental},
				{ID: "node_breakers_model_inference_tripped", Name: "model_inference", Algo: module.Incremental},
				{ID: "node_breakers_accounting_tripped", Name: "accounting", Algo: module.Incremental},
				{ID: "node_breakers_parent_tripped", Name: "parent", Algo: module.Incremental},
			},
		},
	}
)

var clusterHealthCharts = Charts{
	{
		ID:    "cluster_status",
		Title: "Cluster Status",
		Units: "status",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_status",
		Dims: Dims{
			{ID: "cluster_status", Name: "status"},
		},
	},
	{
		ID:    "cluster_number_of_nodes",
		Title: "Cluster Nodes Count",
		Units: "num",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_number_of_nodes",
		Dims: Dims{
			{ID: "cluster_number_of_nodes", Name: "nodes"},
			{ID: "cluster_number_of_data_nodes", Name: "data_nodes"},
		},
	},
	{
		ID:    "cluster_shards_count",
		Title: "Cluster Shards Count",
		Units: "num",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_shards_count",
		Dims: Dims{
			{ID: "cluster_active_primary_shards", Name: "active_primary"},
			{ID: "cluster_active_shards", Name: "active"},
			{ID: "cluster_relocating_shards", Name: "relocating"},
			{ID: "cluster_initializing_shards", Name: "initializing"},
			{ID: "cluster_unassigned_shards", Name: "unassigned"},
			{ID: "cluster_delayed_unassigned_shards", Name: "delayed_unassigned"},
		},
	},
	{
		ID:    "cluster_pending_tasks",
		Title: "Cluster Pending Tasks",
		Units: "num",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_pending_tasks",
		Dims: Dims{
			{ID: "cluster_number_of_pending_tasks", Name: "pending"},
		},
	},
	{
		ID:    "number_of_in_flight_fetch",
		Title: "Unfinished Fetches",
		Units: "num",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.number_of_in_flight_fetch",
		Dims: Dims{
			{ID: "number_of_in_flight_fetch", Name: "in_flight_fetch"},
		},
	},
}

var clusterStatsCharts = Charts{
	{
		ID:    "cluster_indices_count",
		Title: "Cluster Indices Count",
		Units: "num",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_count",
		Dims: Dims{
			{ID: "cluster_indices_count", Name: "indices"},
		},
	},
	{
		ID:    "cluster_indices_shards_count",
		Title: "Cluster Indices Shards Count",
		Units: "num",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_shards_count",
		Type:  module.Area,
		Dims: Dims{
			{ID: "cluster_stats_indices_shards_total", Name: "total"},
			{ID: "cluster_stats_indices_shards_primaries", Name: "primaries"},
			{ID: "cluster_stats_indices_shards_replication", Name: "replication"},
		},
	},
	{
		ID:    "cluster_indices_docs_count",
		Title: "Cluster Indices Docs Count",
		Units: "num",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_docs_count",
		Dims: Dims{
			{ID: "cluster_indices_docs_count", Name: "docs"},
		},
	},
	{
		ID:    "cluster_indices_store_size",
		Title: "Cluster Indices Store Size",
		Units: "bytes",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_store_size",
		Dims: Dims{
			{ID: "cluster_stats_indices_store_size_in_bytes", Name: "size"},
		},
	},
	{
		ID:    "cluster_indices_query_cache",
		Title: "Cluster Indices Query Cache",
		Units: "events/s",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_query_cache",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cluster_stats_indices_query_cache_hit_count", Name: "hit", Algo: module.Incremental},
			{ID: "cluster_stats_indices_query_cache_miss_count", Name: "miss", Algo: module.Incremental},
		},
	},
	{
		ID:    "cluster_nodes_by_role_count",
		Title: "Cluster Nodes By Role Count",
		Units: "num",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_nodes_by_role_count",
		Dims: Dims{
			{ID: "cluster_nodes_count_coordinating_only", Name: "coordinating_only"},
			{ID: "cluster_nodes_count_data", Name: "data"},
			{ID: "cluster_nodes_count_ingest", Name: "ingest"},
			{ID: "cluster_nodes_count_master", Name: "master"},
			{ID: "cluster_nodes_count_ml", Name: "ml"},
			{ID: "cluster_nodes_count_remote_cluster_client", Name: "remote_cluster_client"},
			{ID: "cluster_nodes_count_remote_voting_only", Name: "voting_only"},
		},
	},
}

func (es *Elasticsearch) addIndexToCharts(index string) {

}