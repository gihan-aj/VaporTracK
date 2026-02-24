package grpcapi

import (
	"context"
	"log"

	"vaportrack/engine/internal/domain"
)

// GameTrackerService implements the generated GameTrackerServer interface.
type GameTrackerService struct {
	// gRPC in Go requires embedding this for forward-compatibility
	UnimplementedGameTrackerServer
	repo domain.TrackedGameRepository
}

// a new instance
func NewGameTrackerService(repo domain.TrackedGameRepository) *GameTrackerService {
	return &GameTrackerService{
		repo: repo,
	}
}

// GetTrackedGames is the actual RPC method called by the WinUI client.
func (s *GameTrackerService) GetTrackedGames(ctx context.Context, req *GetTrackedGamesRequest) (*GetTrackedGamesResponse, error){
	log.Printf("Received gRPC request for UserID: %s", req.GetUserId())

	// 1. Fetch data from the database using our Clean Architecture domain interface
	games, err := s.repo.GetActiveGamesByUser(req.GetUserId())
	if err != nil {
		log.Printf("Error fetching games: %v", err)
		return nil, err
	}

	// 2. Map the Domain models to the generated Protobuf models
	var pbGames []*TrackedGame
	for _, g := range games {
		pbGames = append(pbGames, &TrackedGame{
			AppId:       g.AppID,
			Title:       g.Title,
			TargetPrice: g.TargetPrice,
			IsActive:    g.IsActive,
		})
	}

	// 3. Return the response back over the network
	return &GetTrackedGamesResponse{
		Games: pbGames,
	}, nil
}