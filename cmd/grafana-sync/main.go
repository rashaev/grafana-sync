package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rashaev/grafana-sync/internal/config"
	"github.com/rashaev/grafana-sync/internal/grafana"
	"github.com/rashaev/grafana-sync/internal/keycloak"
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
	// Create timer for removing user from Org task
	timerDel := time.NewTimer(cfg.TimeRunDel)

	// Init Keycloak client
	kcloakClient, err := keycloak.New(serverCtx, cfg.KeycloakURL, cfg.KeycloakClientID, cfg.KeycloakClientSecret, cfg.KeycloakUser, cfg.KeycloakPassword)
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
			log.Println("grafana-sync is stopped")
			os.Exit(0)
		case <-timerSync.C:
			log.Println("started process of adding users to Orgs")
			worker.SyncGroups(kcloakClient, grafanaClient, cfg.GroupRegexRO, "Viewer")
			worker.SyncGroups(kcloakClient, grafanaClient, cfg.GroupRegexRW, "Editor")

			// Сбрасываем и перезапускаем таймер
			timerSync.Reset(cfg.TimeRunSync)
		case <-timerDel.C:
			log.Println("started process of removing users from Orgs")
			worker.RemoveFromOrg(kcloakClient, grafanaClient, cfg.GroupRegexRO, cfg.GroupRegexRW)
			timerDel.Reset(cfg.TimeRunDel)
		}
	}
}
