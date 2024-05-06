package main

import (
	"fmt"
	"log"
	"lzhuk/clients/internal/controller"
	"lzhuk/clients/internal/server"
	"lzhuk/clients/internal/utils"
	"lzhuk/clients/pkg/config"
)

// Запуск сервера клиента
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	controller.InitPointApi(cfg)
	// Проверка соединения с сервисом forum-api
	if err := utils.CheckServerConnection(cfg); err != nil {
		fmt.Printf("Соединения с сервисом forum-api не установлено: %v\n", err)
		return
	}
	log.Println("Проверка системы произошла успешно, запущен сервер клиента")
	server.StartServer(cfg)
}
