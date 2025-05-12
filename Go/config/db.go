package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
	mu   sync.Mutex
}

var GlobalDB *DB

func GetExecutableDir() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable path: %v", err)
	}
	return filepath.Join(filepath.Dir(execPath), "load_balancer.db")
}

func NewDB(dbPath string) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS backend_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			backend_id TEXT NOT NULL,
			status TEXT NOT NULL, -- "healthy" hoặc "recovered"
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			request_count INTEGER,
			success_count INTEGER,
			failure_count INTEGER,
			total_latency INTEGER, -- Nanoseconds
			last_latency INTEGER,
			avg_latency INTEGER,
			last_checked DATETIME,
			consecutive_fails INTEGER,
			consecutive_success INTEGER,
			timeout_break INTEGER,
			is_healthy BOOLEAN,
			last_status INTEGER,
			active_connections INTEGER,
			weight INTEGER,
			current_weight INTEGER,
			details TEXT
		)
	`)
	if err != nil {
		db.Close()
		fmt.Printf("failed to create table: %v", err)
	}

	_, err = db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA cache_size = -10000;
		PRAGMA synchronous = OFF;
	`)
	if err != nil {
		db.Close()
		fmt.Printf("failed to optimize database: %v", err)
	}
	GlobalDB = &DB{conn: db}
}

func (d *DB) InsertMetrics(backendID, status string, metrics *Metrics) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.conn.Exec(`
		INSERT INTO backend_metrics (
			backend_id, status, request_count, success_count, failure_count,
			total_latency, last_latency, avg_latency, last_checked,
			consecutive_fails, consecutive_success, timeout_break,
			is_healthy, last_status, active_connections, weight, current_weight, details
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		backendID, status,
		metrics.RequestCount, metrics.SuccessCount, metrics.FailureCount,
		metrics.TotalLatency.Nanoseconds(), metrics.LastLatency.Nanoseconds(), metrics.AvgLatency.Nanoseconds(),
		metrics.LastChecked,
		metrics.ConsecutiveFails, metrics.ConsecutiveSuccess, metrics.TimeoutBreak,
		metrics.IsHealthy, metrics.LastStatus, metrics.ActiveConnections,
		metrics.Weight, metrics.CurrentWeight,
		fmt.Sprintf("Backend %s is %s", backendID, status),
	)
	if err != nil {
		return fmt.Errorf("failed to insert metrics: %v", err)
	}
	log.Printf("Inserted: %s is %s", backendID, status)
	return nil
}

func (d *DB) ReadMetrics() error {
	rows, err := d.conn.Query(`
		SELECT backend_id, status, timestamp, request_count, success_count, failure_count,
		       total_latency, last_latency, avg_latency, last_checked,
		       consecutive_fails, consecutive_success, timeout_break,
		       is_healthy, last_status, active_connections, weight, current_weight, details
		FROM backend_metrics
	`)
	if err != nil {
		return fmt.Errorf("failed to query metrics: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var backendID, status, timestamp, details string
		var requestCount, successCount, failureCount, totalLatency, lastLatency, avgLatency uint64
		var lastChecked time.Time
		var consecutiveFails, consecutiveSuccess, timeoutBreak uint64
		var isHealthy bool
		var lastStatus int
		var activeConnections, weight, currentWeight int64

		err := rows.Scan(
			&backendID, &status, &timestamp, // Thêm timestamp vào đây
			&requestCount, &successCount, &failureCount,
			&totalLatency, &lastLatency, &avgLatency,
			&lastChecked,
			&consecutiveFails, &consecutiveSuccess, &timeoutBreak,
			&isHealthy, &lastStatus, &activeConnections,
			&weight, &currentWeight, &details,
		)
		if err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		fmt.Printf(
			"Backend: %s, Status: %s, Time: %s, RequestCount: %d, SuccessCount: %d, FailureCount: %d, "+
				"TotalLatency: %d, LastLatency: %d, AvgLatency: %d, LastChecked: %s, "+
				"ConsecutiveFails: %d, ConsecutiveSuccess: %d, TimeoutBreak: %d, IsHealthy: %v, "+
				"LastStatus: %d, ActiveConnections: %d, Weight: %d, CurrentWeight: %d, Details: %s\n",
			backendID, status, timestamp, requestCount, successCount, failureCount,
			totalLatency, lastLatency, avgLatency, lastChecked,
			consecutiveFails, consecutiveSuccess, timeoutBreak, isHealthy,
			lastStatus, activeConnections, weight, currentWeight, details,
		)
	}
	return rows.Err()
}
