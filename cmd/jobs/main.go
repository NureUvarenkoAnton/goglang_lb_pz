package main

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/db"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/jobs"
	"NureUvarenkoAnton/unik_go_lb_4/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := db.Connect()
	repo := core.New(db)
	userService := service.NewUserService(repo)
	scheduler, _ := gocron.NewScheduler()
	jobHandler := jobs.NewJobHandler(scheduler)
	jobHandler.RegisterClearUsers(userService)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish

	scheduler.Shutdown()
	db.Close()
}
