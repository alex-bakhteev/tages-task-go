package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"tages-task-go/internal/config"
)

var httpServer *http.Server

// Run запускает HTTP-сервер на указанном хосте и порту
func Run(router http.Handler, wg *sync.WaitGroup) {
	defer wg.Done()
	cfg := config.GetConfig() // Получаем конфигурацию

	// Настройка HTTP-сервера
	address := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	httpServer = &http.Server{
		Addr:    address,
		Handler: router,
	}

	log.Printf("Запуск сервера на %s...\n", address)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// Shutdown корректно завершает работу сервера
func Shutdown(ctx context.Context) error {
	log.Println("Завершение работы сервера...")
	return httpServer.Shutdown(ctx)
}
