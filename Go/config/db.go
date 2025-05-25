package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
			failure_count INTEGER,
			consecutive_fails INTEGER,
			timeout_break INTEGER,
			last_status INTEGER,
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
			backend_id, status, failure_count,
			consecutive_fails, timeout_break,
			last_status, details
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		backendID, status,
		metrics.FailureCount,
		metrics.ConsecutiveFails,
		metrics.TimeoutRate,
		metrics.LastStatus,
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

	log.Printf("Attempting to delete history with ID: %d", id)
	
	result, err := d.conn.Exec("DELETE FROM backend_metrics WHERE id = ?", id)
	if err != nil {
		log.Printf("Database error while deleting history: %v", err)
		return fmt.Errorf("failed to delete error history: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("No rows were affected for ID: %d", id)
		return fmt.Errorf("no history found with ID: %d", id)
	}

	log.Printf("Successfully deleted history with ID: %d", id)
	return nil
}

func (d *DB) DeleteMultipleErrorHistory(ids []int64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(ids) == 0 {
		return fmt.Errorf("no IDs provided for deletion")
	}

	log.Printf("Attempting to delete multiple histories with IDs: %v", ids)

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM backend_metrics WHERE id IN (%s)", strings.Join(placeholders, ","))
	log.Printf("Executing query: %s with args: %v", query, args)

	result, err := d.conn.Exec(query, args...)
	if err != nil {
		log.Printf("Database error while deleting multiple histories: %v", err)
		return fmt.Errorf("failed to delete multiple error history: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("No rows were affected for IDs: %v", ids)
		return fmt.Errorf("no histories found with provided IDs")
	}

	log.Printf("Successfully deleted %d histories", rowsAffected)
	return nil
}
