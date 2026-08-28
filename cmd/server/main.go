package main

import (
	"log"
	"net/http"
	"os"

	"townbasketball/internal/auth"
	"townbasketball/internal/config"
	"townbasketball/internal/guestbook"
	"townbasketball/internal/httpapi"
	"townbasketball/internal/league"
	"townbasketball/internal/media"
	"townbasketball/internal/store"
)

func main() {
	logger := log.New(os.Stdout, "league ", log.LstdFlags)
	cfg := config.FromEnv()
	if err := cfg.Validate(); err != nil {
		logger.Fatal(err)
	}
	db, err := store.Open(cfg.DatabasePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	authService := auth.NewService(db)
	leagueService := league.NewService(db)
	mediaService := media.NewService(db)
	guestbookService := guestbook.NewService(db)
	if err := authService.SeedDefaults(); err != nil {
		logger.Fatal(err)
	}
	if cfg.SeedFixtures {
		if err := leagueService.SeedFixtures(); err != nil {
			logger.Fatal(err)
		}
		if err := mediaService.SeedFixtures(); err != nil {
			logger.Fatal(err)
		}
	}
	server := httpapi.NewServer(leagueService, authService, mediaService, guestbookService, logger)
	logger.Printf("listening on %s", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, server.Handler()); err != nil {
		logger.Fatal(err)
	}
}
