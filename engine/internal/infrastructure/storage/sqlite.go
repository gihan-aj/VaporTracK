package storage

import (
	"database/sql"
	"fmt"
	"vaportrack/engine/internal/domain"

	// Import the sqlite3 driver anonymously so it registers itself
	_ "github.com/mattn/go-sqlite3"
)

// This line does nothing at runtime, but if SQLiteRepository ever fails to 
// implement the interface, the app will refuse to compile.
var _ domain.TrackedGameRepository = (*SQLiteRepository)(nil)

// SQLiteRepository implements domain.TrackedGameRepository
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates the database connection and ensures tables exist.
func NewSQLiteRepository(connectionString string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create Tables (Idempotent: IF NOT EXISTS)
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS TrackedGames (
		AppID TEXT,
		UserID TEXT,
		Title TEXT NOT NULL,
		TargetPrice REAL NOT NULL,
		IsActive BOOLEAN NOT NULL DEFAULT 1,
		PRIMARY KEY (AppID, UserID)
	);

	CREATE TABLE IF NOT EXISTS PriceHistory (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		AppID TEXT NOT NULL,
		SalePrice REAL NOT NULL,
		NormalPrice REAL NOT NULL,
		DateRecorded DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(AppID) REFERENCES TrackedGames(AppID)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	return nil
}

// AddGame saves a new game to the wishlist.
func (r *SQLiteRepository) AddGame(game *domain.TrackedGame) error {
	query := `
		INSERT INTO TrackedGames (AppID, UserID, Title, TargetPrice, IsActive)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(AppID, UserID) DO UPDATE SET
			TargetPrice = excluded.TargetPrice,
			IsActive = excluded.IsActive;
	`
	// We use the "ON CONFLICT" clause to handle "Upserts" (Update if exists, Insert if new)
	_, err := r.db.Exec(query, game.AppID, game.UserID, game.Title, game.TargetPrice, game.IsActive)
	if err != nil {
		return fmt.Errorf("failed to insert game: %w", err)
	}
	return nil
}

// GetActiveGamesByUser returns all games a specific user is tracking.
func (r *SQLiteRepository) GetActiveGamesByUser(userID string) ([]*domain.TrackedGame, error) {
	query := `SELECT AppID, UserID, Title, TargetPrice, IsActive FROM TrackedGames WHERE UserID = ? AND IsActive = 1`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*domain.TrackedGame
	for rows.Next() {
		var g domain.TrackedGame
		if err := rows.Scan(&g.AppID, &g.UserID, &g.Title, &g.TargetPrice, &g.IsActive); err != nil {
			return nil, err
		}
		games = append(games, &g)
	}
	return games, nil
}

// RecordPriceChange logs a new price point into history.
func (r *SQLiteRepository) RecordPriceChange(history *domain.PriceHistory) error {
	query := `INSERT INTO PriceHistory (AppID, SalePrice, NormalPrice, DateRecorded) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, history.AppID, history.SalePrice, history.NormalPrice, history.DateRecorded)
	return err
}