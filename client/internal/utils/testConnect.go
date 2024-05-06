package utils

import (
	"lzhuk/clients/pkg/config"
	"net/http"
	"strconv"
	"time"
)

func CheckServerConnection(cfg config.Config) error {
	serverURL := "http://localhost:" + strconv.Itoa(cfg.PortServer)
	client := http.Client{
		Timeout: time.Second * 2, // Устанавливаем таймаут для проверки соединения
	}
	_, err := client.Get(serverURL)
	return err
}
