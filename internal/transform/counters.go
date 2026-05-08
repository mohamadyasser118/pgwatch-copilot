package transform

// CounterFields maps pgwatch metric names to the jsonb keys
// that are cumulative counters. Only these fields need rate calculation.

var CounterFields = map[string][]string{

	// pg_stat_bgwriter - background writer statistics

	"bgwriter": {
		"buffers_checkpoint", 
		"buffers_clean",      
		"buffers_backend",    
		"buffers_alloc",      
		"checkpoints_timed", 
		"checkpoints_req",    
	},

	// pg_stat_database - per-database statistics

	"db_stats": {
		"xact_commit",    
		"xact_rollback",  
		"blks_read",      
		"blks_hit",     
		"tup_returned",   
		"tup_fetched",   
		"tup_inserted",   
		"tup_updated",    
		"tup_deleted",    
		"temp_files",     
		"temp_bytes",     
		"deadlocks",      
	},

	// pg_stat_user_tables - per-table statistics

	"table_stats": {
		"seq_scan",          
		"seq_tup_read",      
		"idx_scan",          
		"idx_tup_fetch",    
		"n_tup_ins",         
		"n_tup_upd",         
		"n_tup_del",        
		"n_tup_hot_upd",   
		"vacuum_count",      
		"autovacuum_count",  
		"analyze_count",    
		"autoanalyze_count", 
	},

	// pg_stat_statements - query statistics

	"statements": {
		"calls",           
		"total_exec_time",  
		"rows",            
		"shared_blks_hit",  
		"shared_blks_read", 
		"temp_blks_read",   
		"temp_blks_written",
	},

	// WAL (Write-Ahead Log) statistics

	"wal": {
		"wal_records", 
		"wal_fpi",     
		"wal_bytes",   
	},

	// pg_stat_user_indexes - index usage

	"index_stats": {
		"idx_scan",      
		"idx_tup_read", 
		"idx_tup_fetch", 
	},
}