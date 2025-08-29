package main

import (
	"log"
	wb_task_L0 "wb-task-L0"
	"wb-task-L0/internal/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(wb_task_L0.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}

}
