package main

import (
	"context"
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

	// Загружаем конфиги
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Инициализация БД
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

	// Репозитории и кэш
	repos := repository.NewRepository(db)
	orderCache := cache.NewCache()

	// Загружаем кэш из БД
	orders, err := repos.Order.GetAll()
	if err != nil {
		logrus.Fatalf("failed to load orders: %v", err)
	}
	orderCache.LoadFromDB(orders)
	logrus.Printf("Cache initialized with %d orders", orderCache.Len())

	// Сервисы и HTTP хендлеры
	services := service.NewService(repos, orderCache)
	handlers := handler.NewHandler(services)

	// HTTP сервер
	srv := new(wb_task_L0.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("HTTP server started")

	// Kafka consumer
	brokers := []string{"localhost:9092"} // замените на реальные адреса
	topic := "orders"
	groupID := "order-consumers"

	consumer := kafka.NewConsumer(brokers, topic, groupID, repos.Order, orderCache)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)
	logrus.Print("Kafka consumer started")

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("Shutting down application...")

	// Остановка Kafka consumer
	cancel()
	if err := consumer.Close(); err != nil {
		logrus.Errorf("error closing Kafka consumer: %s", err.Error())
	}

	// Остановка HTTP сервера
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	// Закрытие соединения с БД
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("failed to get sql.DB from gorm: %s", err.Error())
	} else {
		if err := sqlDB.Close(); err != nil {
			logrus.Errorf("error occured on db connection close: %s", err.Error())
		}
	}

	// Миграции (на всякий случай)
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
