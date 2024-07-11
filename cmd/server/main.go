package main

import (
	"fmt"
	"github.com/instinctG/statistics_service/internal/db"
	"github.com/instinctG/statistics_service/internal/statistics"
	transportHttp "github.com/instinctG/statistics_service/internal/transport/http"
)

// Run инициализирует и запускает приложение, устанавливает соединение с базой данных,
// выполняет миграции базы данных, создает сервис статистики и запускает HTTP сервер.
func Run() error {
	fmt.Println("starting up our application")

	database, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}

	if err := database.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	statisticService := statistics.NewService(database)

	httpHandler := transportHttp.NewHandler(statisticService)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

// main является точкой входа в приложение, вызывает функцию Run и обрабатывает возможные ошибки.
func main() {
	fmt.Println("Running statistics service")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
