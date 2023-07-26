package main

import (
	"log"

	todo "github.com/WalkerChel/TO-DO"
	"github.com/WalkerChel/TO-DO/pkg/handler"
	"github.com/WalkerChel/TO-DO/pkg/repository"
	"github.com/WalkerChel/TO-DO/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)
	srv := new(todo.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
