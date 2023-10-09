package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rashaev/grafana-sync/internal/avanpost"
	"github.com/rashaev/grafana-sync/internal/config"
	"github.com/rashaev/grafana-sync/internal/grafana"
	"github.com/rashaev/grafana-sync/internal/worker"
)

func main() {

	serverCtx := context.Background()

	// Read Config from env
	cfg := config.InitConfig()

	// Create channel for shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Create timer for sync user task
	timerSync := time.NewTimer(cfg.TimeRunSync)
	// Create timer for delete task
	timerDel := time.NewTimer(10 * time.Second)

	// Init Avanpost client
	avpClient, err := avanpost.New(serverCtx, cfg.AvanpostURL, cfg.AvanpostClientID, cfg.AvanpostClientSecret, cfg.AvanpostUser, cfg.AvanpostPassword)
	if err != nil {
		log.Fatalln(err)
	}

	// Init Grafana client
	grafanaClient, err := grafana.New(cfg.GrafanaURL, grafana.NewConfig(cfg.GrafanaUser, cfg.GrafanaPassword))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case <-stopChan:
			// Прерывание программы
			log.Println("grafana-sync is stopped")
			return
		case <-timerSync.C:
			// Действия, выполняемые по истечении таймера
			worker.SyncGroups(avpClient, grafanaClient, cfg.GroupRegexRO)

			// Сбрасываем и перезапускаем таймер
			timerSync.Reset(cfg.TimeRunSync)
		case <-timerDel.C:
			worker.DeleteGrafanaUser(grafanaClient)
			timerDel.Reset(10 * time.Second)
		}
	}
}
