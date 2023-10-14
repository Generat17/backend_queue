package main

import (
	"context"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"server"
	"server/pkg/handler"
	"server/pkg/repository"
	"server/pkg/service"
	"syscall"
	"time"
)

// @title       Digital Queue API
// @version     1.0
// @description API Server for digital queue

// @contact.name   Timur Aliev
// @contact.url    https://t.me/Aliev_Timur_M
// @contact.email  alievtm@gmail.com

// @host     localhost:8000
// @BasePath /api/v1

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	var checkUpdate [1]int
	var oldCheckUpdate int = 0

	repos := repository.NewRepository(db)
	services := service.NewService(repos, &checkUpdate)
	handlers := handler.NewHandler(services)

	// инициализируем переменную, в которой будем хранить текущий статус сотрудников
	//employeeStatus, _ := services.Employee.GetEmployeeStatusList()

	srv := new(server.Server)

	// new code (test SSE)
	// Create SSE server
	s := sse.NewServer(nil)
	defer s.Shutdown()

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes(s)); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	go func() {
		for {
			services.Queue.CheckLongWait()
			time.Sleep(60 * time.Second)
		}
	}()

	// Send messages
	go func() {
		for {
			if oldCheckUpdate != checkUpdate[0] {
				oldCheckUpdate = checkUpdate[0]
				s.SendMessage("/events/channel", sse.SimpleMessage("updated"))
				logrus.Println("updated")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	logrus.Print("Server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("configs")
	return viper.ReadInConfig()
}

// Команда для запуска линтера:  golangci-lint run
// Команда для init swagger:
