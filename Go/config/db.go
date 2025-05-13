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

type ErrorBackend struct {
	ID               int64
	Backend          string
	Time             time.Time
	Status           string
	Details          string
	ConsecutiveFails int
	TimeoutBreak     int
	LastStatus       int
	FailureCount     int
}

func (d *DB) ReadMetrics() []ErrorBackend {
	rows, err := d.conn.Query(`
		SELECT id, backend_id, status, timestamp, failure_count,
		    consecutive_fails, timeout_break,
		    last_status, details
		FROM backend_metrics
	`)
	if err != nil {
		fmt.Println("failed to query metrics: ", err)
	}
	defer rows.Close()

	metricsMap := []ErrorBackend{}

	for rows.Next() {
		var e ErrorBackend
		err := rows.Scan(
			&e.ID,
			&e.Backend, &e.Status, &e.Time, &e.FailureCount,
			&e.ConsecutiveFails, &e.TimeoutBreak, &e.LastStatus, &e.Details,
		)
		if err != nil {
			fmt.Println("failed to scan error row: ", err)
		}
		metricsMap = append(metricsMap, e)
	}
	return metricsMap
}

func (d *DB) DeleteErrorHistory(id int64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.conn.Exec(`
		DELETE FROM backend_metrics 
		WHERE id = ?
	`, id)
	
	if err != nil {
		return fmt.Errorf("failed to delete error history: %v", err)
	}
	return nil
}
