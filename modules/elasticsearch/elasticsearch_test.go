package elasticsearch

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v790SingleNodesLocalStats, _ = ioutil.ReadFile("testdata/v790_single_nodes_local_stats.json")
	v790SingleClusterHealth, _   = ioutil.ReadFile("testdata/v790_single_cluster_health.json")
	v790SingleClusterStats, _    = ioutil.ReadFile("testdata/v790_single_cluster_stats.json")
	v790SingleCatIndicesStats, _ = ioutil.ReadFile("testdata/v790_single_cat_indices_stats.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"v790SingleNodesLocalStats": v790SingleNodesLocalStats,
		"v790SingleClusterHealth":   v790SingleClusterHealth,
		"v790SingleClusterStats":    v790SingleClusterStats,
		"v790SingleCatIndicesStats": v790SingleCatIndicesStats,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestElasticsearch_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default": {
			config: New().Config,
		},
		"URL not set": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{UserURL: ""},
				}},
		},
		"invalid TLSCA": {
			wantFail: true,
			config: Config{
				HTTP: web.HTTP{
					Client: web.Client{ClientTLSConfig: web.ClientTLSConfig{TLSCA: "testdata/tls"}},
				}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es := New()
			es.Config = test.config

			if test.wantFail {
				assert.False(t, es.Init())
			} else {
				assert.True(t, es.Init())
			}
		})
	}
}

func TestElasticsearch_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(*testing.T) (prom *Elasticsearch, cleanup func())
		wantFail bool
	}{
		"valid data":         {prepare: prepareElasticsearchValidData},
		"invalid data":       {prepare: prepareElasticsearchInvalidData, wantFail: true},
		"404":                {prepare: prepareElasticsearch404, wantFail: true},
		"connection refused": {prepare: prepareElasticsearchConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prom, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, prom.Check())
			} else {
				assert.True(t, prom.Check())
			}
		})
	}
}

func TestElasticsearch_Charts(t *testing.T) {
	assert.Nil(t, New().Charts())
}

func TestElasticsearch_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestElasticsearch_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *Elasticsearch
		wantCollected map[string]int64
	}{
		"v790: all stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = true
				es.DoClusterHealth = true
				es.DoClusterStats = true
				es.DoIndicesStats = true
				return es
			},
			wantCollected: map[string]int64{
				"cluster_active_primary_shards":                                0,
				"cluster_active_shards":                                        0,
				"cluster_active_shards_percent_as_number":                      100,
				"cluster_delayed_unassigned_shards":                            0,
				"cluster_indices_count":                                        0,
				"cluster_indices_docs_count":                                   0,
				"cluster_indices_query_cache_hit_count":                        0,
				"cluster_indices_query_cache_miss_count":                       0,
				"cluster_indices_shards_primaries":                             0,
				"cluster_indices_shards_replication":                           0,
				"cluster_indices_shards_total":                                 0,
				"cluster_indices_store_size_in_bytes":                          0,
				"cluster_initializing_shards":                                  0,
				"cluster_nodes_count_coordinating_only":                        0,
				"cluster_nodes_count_data":                                     1,
				"cluster_nodes_count_ingest":                                   1,
				"cluster_nodes_count_master":                                   1,
				"cluster_nodes_count_ml":                                       1,
				"cluster_nodes_count_remote_cluster_client":                    1,
				"cluster_nodes_count_total":                                    1,
				"cluster_nodes_count_transform":                                1,
				"cluster_nodes_count_voting_only":                              0,
				"cluster_number_of_data_nodes":                                 1,
				"cluster_number_of_in_flight_fetch":                            0,
				"cluster_number_of_nodes":                                      1,
				"cluster_number_of_pending_tasks":                              0,
				"cluster_relocating_shards":                                    0,
				"cluster_status":                                               0,
				"cluster_unassigned_shards":                                    0,
				"node_breakers_accounting_tripped":                             0,
				"node_breakers_fielddata_tripped":                              0,
				"node_breakers_in_flight_requests_tripped":                     0,
				"node_breakers_model_inference_tripped":                        0,
				"node_breakers_parent_tripped":                                 0,
				"node_breakers_requests_tripped":                               0,
				"node_http_current_open":                                       3,
				"node_indices_fielddata_evictions":                             0,
				"node_indices_fielddata_memory_size_in_bytes":                  0,
				"node_indices_flush_total":                                     0,
				"node_indices_flush_total_time_in_millis":                      0,
				"node_indices_indexing_index_current":                          0,
				"node_indices_indexing_index_time_in_millis":                   0,
				"node_indices_indexing_index_total":                            0,
				"node_indices_refresh_total":                                   0,
				"node_indices_refresh_total_time_in_millis":                    0,
				"node_indices_search_fetch_current":                            0,
				"node_indices_search_fetch_time_in_millis":                     0,
				"node_indices_search_fetch_total":                              0,
				"node_indices_search_query_current":                            0,
				"node_indices_search_query_time_in_millis":                     0,
				"node_indices_search_query_total":                              0,
				"node_indices_segments_count":                                  0,
				"node_indices_segments_doc_values_memory_in_bytes":             0,
				"node_indices_segments_fixed_bit_set_memory_in_bytes":          0,
				"node_indices_segments_index_writer_memory_in_bytes":           0,
				"node_indices_segments_memory_in_bytes":                        0,
				"node_indices_segments_norms_memory_in_bytes":                  0,
				"node_indices_segments_points_memory_in_bytes":                 0,
				"node_indices_segments_stored_fields_memory_in_bytes":          0,
				"node_indices_segments_term_vectors_memory_in_bytes":           0,
				"node_indices_segments_terms_memory_in_bytes":                  0,
				"node_indices_segments_version_map_memory_in_bytes":            0,
				"node_indices_stats_my-index-000001_index_docs_count":          0,
				"node_indices_stats_my-index-000001_index_health":              1,
				"node_indices_stats_my-index-000001_index_shards_count":        1,
				"node_indices_stats_my-index-000001_index_store_size_in_bytes": 208,
				"node_indices_stats_my-index-000002_index_docs_count":          0,
				"node_indices_stats_my-index-000002_index_health":              1,
				"node_indices_stats_my-index-000002_index_shards_count":        1,
				"node_indices_stats_my-index-000002_index_store_size_in_bytes": 208,
				"node_indices_stats_my-index-000003_index_docs_count":          0,
				"node_indices_stats_my-index-000003_index_health":              1,
				"node_indices_stats_my-index-000003_index_shards_count":        1,
				"node_indices_stats_my-index-000003_index_store_size_in_bytes": 208,
				"node_indices_translog_operations":                             0,
				"node_indices_translog_size_in_bytes":                          0,
				"node_indices_translog_uncommitted_operations":                 0,
				"node_indices_translog_uncommitted_size_in_bytes":              0,
				"node_jvm_buffer_pool_direct_count":                            0,
				"node_jvm_buffer_pool_direct_total_capacity_in_bytes":          0,
				"node_jvm_buffer_pool_direct_used_in_bytes":                    0,
				"node_jvm_buffer_pool_mapped_count":                            0,
				"node_jvm_buffer_pool_mapped_total_capacity_in_bytes":          0,
				"node_jvm_buffer_pool_mapped_used_in_bytes":                    0,
				"node_jvm_gc_collectors_old_collection_count":                  0,
				"node_jvm_gc_collectors_old_collection_time_in_millis":         0,
				"node_jvm_gc_collectors_young_collection_count":                16,
				"node_jvm_gc_collectors_young_collection_time_in_millis":       184,
				"node_jvm_mem_heap_committed_in_bytes":                         1073741824,
				"node_jvm_mem_heap_used_in_bytes":                              363166720,
				"node_jvm_mem_heap_used_percent":                               33,
				"node_process_max_file_descriptors":                            1048576,
				"node_process_open_file_descriptors":                           258,
				"node_thread_pool_search_queue":                                0,
				"node_thread_pool_search_rejected":                             0,
				"node_thread_pool_write_queue":                                 0,
				"node_thread_pool_write_rejected":                              0,
				"node_transport_rx_count":                                      0,
				"node_transport_rx_size_in_bytes":                              0,
				"node_transport_tx_count":                                      0,
				"node_transport_tx_size_in_bytes":                              0,
			},
		},
		"v790: only node_stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = true
				es.DoClusterHealth = false
				es.DoClusterStats = false
				es.DoIndicesStats = false
				return es
			},
			wantCollected: map[string]int64{
				"node_breakers_accounting_tripped":                       0,
				"node_breakers_fielddata_tripped":                        0,
				"node_breakers_in_flight_requests_tripped":               0,
				"node_breakers_model_inference_tripped":                  0,
				"node_breakers_parent_tripped":                           0,
				"node_breakers_requests_tripped":                         0,
				"node_http_current_open":                                 3,
				"node_indices_fielddata_evictions":                       0,
				"node_indices_fielddata_memory_size_in_bytes":            0,
				"node_indices_flush_total":                               0,
				"node_indices_flush_total_time_in_millis":                0,
				"node_indices_indexing_index_current":                    0,
				"node_indices_indexing_index_time_in_millis":             0,
				"node_indices_indexing_index_total":                      0,
				"node_indices_refresh_total":                             0,
				"node_indices_refresh_total_time_in_millis":              0,
				"node_indices_search_fetch_current":                      0,
				"node_indices_search_fetch_time_in_millis":               0,
				"node_indices_search_fetch_total":                        0,
				"node_indices_search_query_current":                      0,
				"node_indices_search_query_time_in_millis":               0,
				"node_indices_search_query_total":                        0,
				"node_indices_segments_count":                            0,
				"node_indices_segments_doc_values_memory_in_bytes":       0,
				"node_indices_segments_fixed_bit_set_memory_in_bytes":    0,
				"node_indices_segments_index_writer_memory_in_bytes":     0,
				"node_indices_segments_memory_in_bytes":                  0,
				"node_indices_segments_norms_memory_in_bytes":            0,
				"node_indices_segments_points_memory_in_bytes":           0,
				"node_indices_segments_stored_fields_memory_in_bytes":    0,
				"node_indices_segments_term_vectors_memory_in_bytes":     0,
				"node_indices_segments_terms_memory_in_bytes":            0,
				"node_indices_segments_version_map_memory_in_bytes":      0,
				"node_indices_translog_operations":                       0,
				"node_indices_translog_size_in_bytes":                    0,
				"node_indices_translog_uncommitted_operations":           0,
				"node_indices_translog_uncommitted_size_in_bytes":        0,
				"node_jvm_buffer_pool_direct_count":                      0,
				"node_jvm_buffer_pool_direct_total_capacity_in_bytes":    0,
				"node_jvm_buffer_pool_direct_used_in_bytes":              0,
				"node_jvm_buffer_pool_mapped_count":                      0,
				"node_jvm_buffer_pool_mapped_total_capacity_in_bytes":    0,
				"node_jvm_buffer_pool_mapped_used_in_bytes":              0,
				"node_jvm_gc_collectors_old_collection_count":            0,
				"node_jvm_gc_collectors_old_collection_time_in_millis":   0,
				"node_jvm_gc_collectors_young_collection_count":          16,
				"node_jvm_gc_collectors_young_collection_time_in_millis": 184,
				"node_jvm_mem_heap_committed_in_bytes":                   1073741824,
				"node_jvm_mem_heap_used_in_bytes":                        363166720,
				"node_jvm_mem_heap_used_percent":                         33,
				"node_process_max_file_descriptors":                      1048576,
				"node_process_open_file_descriptors":                     258,
				"node_thread_pool_search_queue":                          0,
				"node_thread_pool_search_rejected":                       0,
				"node_thread_pool_write_queue":                           0,
				"node_thread_pool_write_rejected":                        0,
				"node_transport_rx_count":                                0,
				"node_transport_rx_size_in_bytes":                        0,
				"node_transport_tx_count":                                0,
				"node_transport_tx_size_in_bytes":                        0,
			},
		},
		"v790: only cluster_health": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = false
				es.DoClusterHealth = true
				es.DoClusterStats = false
				es.DoIndicesStats = false
				return es
			},
			wantCollected: map[string]int64{
				"cluster_active_primary_shards":           0,
				"cluster_active_shards":                   0,
				"cluster_active_shards_percent_as_number": 100,
				"cluster_delayed_unassigned_shards":       0,
				"cluster_initializing_shards":             0,
				"cluster_number_of_data_nodes":            1,
				"cluster_number_of_in_flight_fetch":       0,
				"cluster_number_of_nodes":                 1,
				"cluster_number_of_pending_tasks":         0,
				"cluster_relocating_shards":               0,
				"cluster_status":                          0,
				"cluster_unassigned_shards":               0,
			},
		},
		"v790: only cluster_stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = false
				es.DoClusterHealth = false
				es.DoClusterStats = true
				es.DoIndicesStats = false
				return es
			},
			wantCollected: map[string]int64{
				"cluster_indices_count":                     0,
				"cluster_indices_docs_count":                0,
				"cluster_indices_query_cache_hit_count":     0,
				"cluster_indices_query_cache_miss_count":    0,
				"cluster_indices_shards_primaries":          0,
				"cluster_indices_shards_replication":        0,
				"cluster_indices_shards_total":              0,
				"cluster_indices_store_size_in_bytes":       0,
				"cluster_nodes_count_coordinating_only":     0,
				"cluster_nodes_count_data":                  1,
				"cluster_nodes_count_ingest":                1,
				"cluster_nodes_count_master":                1,
				"cluster_nodes_count_ml":                    1,
				"cluster_nodes_count_remote_cluster_client": 1,
				"cluster_nodes_count_total":                 1,
				"cluster_nodes_count_transform":             1,
				"cluster_nodes_count_voting_only":           0,
			},
		},
		"v790: only indices_stats": {
			prepare: func() *Elasticsearch {
				es := New()
				es.DoNodeStats = false
				es.DoClusterHealth = false
				es.DoClusterStats = false
				es.DoIndicesStats = true
				return es
			},
			wantCollected: map[string]int64{
				"node_indices_stats_my-index-000001_index_docs_count":          0,
				"node_indices_stats_my-index-000001_index_health":              1,
				"node_indices_stats_my-index-000001_index_shards_count":        1,
				"node_indices_stats_my-index-000001_index_store_size_in_bytes": 208,
				"node_indices_stats_my-index-000002_index_docs_count":          0,
				"node_indices_stats_my-index-000002_index_health":              1,
				"node_indices_stats_my-index-000002_index_shards_count":        1,
				"node_indices_stats_my-index-000002_index_store_size_in_bytes": 208,
				"node_indices_stats_my-index-000003_index_docs_count":          0,
				"node_indices_stats_my-index-000003_index_health":              1,
				"node_indices_stats_my-index-000003_index_shards_count":        1,
				"node_indices_stats_my-index-000003_index_store_size_in_bytes": 208,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es, cleanup := prepareElasticsearch(t, test.prepare)
			defer cleanup()

			var collected map[string]int64
			for i := 0; i < 10; i++ {
				collected = es.Collect()
			}

			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, es, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, es *Elasticsearch, collected map[string]int64) {
	for _, chart := range *es.Charts() {
		if chart.Obsolete {
			continue
		}
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareElasticsearch(t *testing.T, createES func() *Elasticsearch) (es *Elasticsearch, cleanup func()) {
	t.Helper()
	srv := prepareElasticSearchEndpoint()

	es = createES()
	es.UserURL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

func prepareElasticsearchValidData(t *testing.T) (es *Elasticsearch, cleanup func()) {
	return prepareElasticsearch(t, New)
}

func prepareElasticsearchInvalidData(t *testing.T) (*Elasticsearch, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))
	es := New()
	es.UserURL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

func prepareElasticsearch404(t *testing.T) (*Elasticsearch, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	es := New()
	es.UserURL = srv.URL
	require.True(t, es.Init())

	return es, srv.Close
}

func prepareElasticsearchConnectionRefused(t *testing.T) (*Elasticsearch, func()) {
	t.Helper()
	es := New()
	es.UserURL = "http://127.0.0.1:38001/metrics"
	require.True(t, es.Init())

	return es, func() {}
}

func prepareElasticSearchEndpoint() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case urlPathLocalNodeStats:
				_, _ = w.Write(v790SingleNodesLocalStats)
			case urlPathClusterHealth:
				_, _ = w.Write(v790SingleClusterHealth)
			case urlPathClusterStats:
				_, _ = w.Write(v790SingleClusterStats)
			case urlPathIndicesStats:
				_, _ = w.Write(v790SingleCatIndicesStats)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}
