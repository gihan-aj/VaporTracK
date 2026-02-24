package main

import (
	"log"
	"vaportrack/engine/internal/domain"
	"vaportrack/engine/internal/infrastructure/storage"
)

func main() {
	log.Println("Booting up VaporTrack Engine...")

	// 1. Initialize Infrastructure (SQLite)
	// This will create 'vaportrack.db' in the folder where you run the command
	repo, err := storage.NewSQLiteRepository("vaportrack.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection established and tables created!")
	
	// 2. Test the Db: Add Cyberpunk 2077
	log.Println("Seeding database with Cyberpunk 2077...")
	cyberpunk := &domain.TrackedGame{
		AppID: 			"1091500",
		UserID: 		"user_1",
		Title: 			"Cyberpunk 2077",
		TargetPrice:	29.99,
		IsActive: 		true,
	}
	
	if err := repo.AddGame(cyberpunk); err != nil {
		log.Printf("Failed to add game: %v", err)
	} else {
		log.Printf("Successfully added Cyberpunk 2077 to the wishlist!")
	}

	// 3. Verify: Read it back
	games, err := repo.GetActiveGamesByUser("user_1")
	if err != nil {
		log.Printf("Failed to fetch games: %v", err)
	}

	log.Printf("User_1 is currently tracking %d games", len(games))
	for _, g := range games {
		log.Printf(" - %s (Target: $%.2f)", g.Title, g.TargetPrice)
	}

	// Initialize Domain Services / Use Cases
	// TODO: Create the tracking service that uses the repository

	// Start the Engine (gRPC Server / Background Polling)
	// TODO: Start goroutines for polling CheapShark

	log.Println("Engine is running. Press CTRL+C to exit.")
	
	// Temporary block to keep the application running until we add real background workers
	select {} 
}