package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/WalkerChel/TO-DO/configs"
	"github.com/WalkerChel/TO-DO/internal/handler"
	"github.com/WalkerChel/TO-DO/internal/repo"
	"github.com/WalkerChel/TO-DO/internal/service"
	"github.com/WalkerChel/TO-DO/pkg/httpserver"
	"github.com/WalkerChel/TO-DO/pkg/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Run(ConfigPath string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cnf, err := configs.New(ConfigPath)
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
	db, err := postgres.NewPostgresDB(postgres.ConfigPG{
		Host:     cnf.PG.Host,
		Port:     cnf.PG.Port,
		Username: cnf.PG.Username,
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   cnf.PG.DBname,
		SSLMode:  cnf.PG.SslMode,
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repo.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := new(httpserver.Server)

	go func() {
		if err := srv.Run(cnf.HTTP.Port, handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Printf("TO-DO app is running on port: %s", cnf.HTTP.Port)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("TO-DO app shutting down")

	if err := srv.ShutDown(context.Background()); err != nil {
		logrus.Errorf("error while shuttung down TO-DO app: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error while closing database connection: %s", err.Error())
	}

}
