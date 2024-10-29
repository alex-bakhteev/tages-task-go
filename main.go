package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tages-task-go/cmd/server"
	"time"
)

func main() {
	// Инициализация приложения
	router, err := server.Initialize()
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}
	defer server.DbPool.Close()

	// Используем WaitGroup для управления горутинами
	var wg sync.WaitGroup
	wg.Add(1)

	// Запуск сервера в отдельной горутине
	go server.Run(router, &wg)

	// Перехват сигнала для корректного завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Получен сигнал для завершения")

	// Инициализация контекста для shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершение работы HTTP-сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	// Ожидание завершения горутины сервера
	wg.Wait()
	log.Println("Приложение завершено.")
}
