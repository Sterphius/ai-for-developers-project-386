package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sterphius/ai-for-developers-project-386/internal/config"
	httpapi "github.com/Sterphius/ai-for-developers-project-386/internal/httpapi"
	"github.com/Sterphius/ai-for-developers-project-386/internal/repository/postgres"
	"github.com/Sterphius/ai-for-developers-project-386/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create db pool: %v", err)
	}
	defer pool.Close()

	repo := postgres.New(pool)
	appService := service.New(repo, service.Settings{
		OwnerID:             cfg.OwnerID,
		Location:            cfg.Location,
		WindowDays:          cfg.WindowDays,
		WorkingDayStartHour: cfg.WorkingDayStartHour,
		WorkingDayEndHour:   cfg.WorkingDayEndHour,
		SlotStepMinutes:     cfg.SlotStepMinutes,
	})

	router := httpapi.NewRouter(appService)

	server := httpapi.NewServer(router, cfg.ListenAddr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
