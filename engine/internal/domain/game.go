package domain

import "time"

// User's wishlist items
type TrackedGame struct {
	AppID		string
	UserID		string
	Title		string
	TargetPrice	float64
	IsActive	bool
}

// Single ledger entry for a games price at a specific time
type PriceHistory struct {
	ID				int64
	AppID 			string
	SalePrice		float64
	NormalPrice		float64
	DateRecorded	time.Time
}

// TrackedGameRepository defines exactly how the rest of the app can interact with storage.
// Notice this is an interface. The domain doesn't care IF this is SQLite, Postgres, or Memory.
type TrackedGameRepository interface {
	AddGame(game *TrackedGame) error
	GetActiveGamesByUser(userID string) ([]*TrackedGame, error)
	RecordPriceChange(history *PriceHistory) error
}