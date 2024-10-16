package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amit8889/golangCRUDApi/internal/config"
	student "github.com/amit8889/golangCRUDApi/internal/http/handlers"
	sqllite "github.com/amit8889/golangCRUDApi/internal/storage/sqlite"
)

func main() {
	fmt.Println("========Welcome to sudents-api======")
	//load config
	cfg := config.MustLoad()
	//db setup
	storage, err := sqllite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initalize", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//router setup
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/getStudentById/{id}", student.GetStudentByID(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}
	fmt.Printf("Server started %s", cfg.HttpServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGALRM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("==>Failed to start server")
		}
	}()
	<-done
	slog.Info("Server stopped")
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")

}
