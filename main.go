package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"young-astrologer-Nastenka/internal/api"
	"young-astrologer-Nastenka/internal/app"
	"young-astrologer-Nastenka/internal/repository"
	"young-astrologer-Nastenka/internal/service"
)

func main() {
	cfg, err := app.NewConfig()
	if err != nil {
		log.Panicf("get config: %s", err)
	}

	db, err := app.ConnectToPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Panicf("connect to postgres: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	repo := repository.NewRepository(db)
	s := service.NewService(cfg.NASAAPI, repo)
	h := api.NewHandler(s)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.UpdateAPODJob(ctx)

	srv := api.NewServer(cfg.HTTPPort, h)

	go func() {
		err = srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-c

	downCtx, _ := context.WithTimeout(ctx, time.Second)
	err = srv.Shutdown(downCtx)
	if err != nil {
		log.Panic(err)
	}
}
