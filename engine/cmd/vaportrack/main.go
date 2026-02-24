package main

import (
	"log"
	"net"

	"vaportrack/engine/internal/domain"
	"vaportrack/engine/internal/infrastructure/storage"
	grpcapi "vaportrack/engine/internal/infrastructure/grpc" // Alias the import

	"google.golang.org/grpc"
)

func main() {
	log.Println("Booting up VaporTrack Engine...")

	// 1. Initialize Infrastructure (SQLite)
	repo, err := storage.NewSQLiteRepository("vaportrack.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection established!")

	// (Optional) Seed the database with Cyberpunk 2077 if it's not already there
	cyberpunk := &domain.TrackedGame{
		AppID:       "1091500",
		UserID:      "user_1",
		Title:       "Cyberpunk 2077",
		TargetPrice: 29.99,
		IsActive:    true,
	}
	_ = repo.AddGame(cyberpunk) // Ignoring error for brevity since it uses UPSERT

	// 2. Initialize the gRPC Service
	trackerService := grpcapi.NewGameTrackerService(repo)

	// 3. Set up the TCP Listener on port 50051 (Standard gRPC port)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// 4. Create and start the gRPC Server
	grpcServer := grpc.NewServer()
	grpcapi.RegisterGameTrackerServer(grpcServer, trackerService)

	log.Println("⚡ gRPC Server is actively listening on port 50051...")
	
	// Serve() is a blocking call, so we no longer need the 'select {}' trick
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}