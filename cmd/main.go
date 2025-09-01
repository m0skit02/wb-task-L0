package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
	"wb-task-L0/pkg/cache"
	"wb-task-L0/pkg/kafka"
	"wb-task-L0/pkg/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	wb_task_L0 "wb-task-L0"
	"wb-task-L0/pkg/handler"
	"wb-task-L0/pkg/repository"
	"wb-task-L0/pkg/service"
)

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

	repos := repository.NewRepository(db)
	orderCache := cache.NewCache()

	orders, err := repos.Order.GetAll()
	if err != nil {
		logrus.Fatalf("failed to load orders: %v", err)
	}
	orderCache.LoadFromDB(orders)
	logrus.Printf("Cache initialized with %d orders", orderCache.Len())

	services := service.NewService(repos, orderCache)
	handlers := handler.NewHandler(services)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.Static("/static", "./web")
	router.GET("/", func(c *gin.Context) {
		c.File("./web/front.html")
	})

	apiRoutes := handlers.InitRoutes()
	router.Any("/api/*any", gin.WrapH(apiRoutes))

	srv := &wb_task_L0.Server{}
	go func() {
		if err := srv.Run(viper.GetString("port"), router); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("HTTP server started")

	brokerEnv := os.Getenv("KAFKA_BROKER")
	topicEnv := os.Getenv("KAFKA_TOPIC")
	if brokerEnv == "" || topicEnv == "" {
		logrus.Fatal("KAFKA_BROKER or KAFKA_TOPIC is not set in environment")
	}
	brokers := []string{brokerEnv}
	topic := topicEnv
	groupID := "order-consumers"

	consumer := kafka.NewConsumer(brokers, topic, groupID, repos.Order, orderCache)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)
	logrus.Print("Kafka consumer started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("Shutting down application...")

	cancel()
	if err := consumer.Close(); err != nil {
		logrus.Errorf("error closing Kafka consumer: %s", err.Error())
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("failed to get sql.DB from gorm: %s", err.Error())
	} else {
		if err := sqlDB.Close(); err != nil {
			logrus.Errorf("error occured on db connection close: %s", err.Error())
		}
	}

	if err := db.AutoMigrate(&models.Order{}, &models.Delivery{}, &models.Payment{}, &models.Item{}); err != nil {
		logrus.Fatalf("failed to migrate: %s", err.Error())
	}

	logrus.Print("Application shutdown complete")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
