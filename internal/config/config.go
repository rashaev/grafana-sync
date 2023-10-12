package config

import (
	"log"
	"net/url"
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	GrafanaURL           url.URL       `env:"GRAFANA_SYNC_GRAFANA_URL,notEmpty"`
	GrafanaUser          string        `env:"GRAFANA_SYNC_GRAFANA_USER,notEmpty"`
	GrafanaPassword      string        `env:"GRAFANA_SYNC_GRAFANA_PASSWORD,notEmpty"`
	KeycloakURL          url.URL       `env:"GRAFANA_SYNC_KEYCLOAK_URL,notEmpty"`
	KeycloakUser         string        `env:"GRAFANA_SYNC_KEYCLOAK_USER,notEmpty"`
	KeycloakPassword     string        `env:"GRAFANA_SYNC_KEYCLOAK_PASSWORD,notEmpty"`
	KeycloakClientID     string        `env:"GRAFANA_SYNC_KEYCLOAK_CLIENT_ID,notEmpty"`
	KeycloakClientSecret string        `env:"GRAFANA_SYNC_KEYCLOAK_CLIENT_SECRET,notEmpty"`
	GroupRegexRO         string        `env:"GRAFANA_SYNC_GROUP_REGEX_RO,notEmpty"`
	GroupRegexRW         string        `env:"GRAFANA_SYNC_GROUP_REGEX_RW,notEmpty"`
	TimeRunSync          time.Duration `env:"GRAFANA_SYNC_TIME_RUN_SYNC,notEmpty" envDefault:"20m"`
	TimeRunDel           time.Duration `env:"GRAFANA_SYNC_TIME_RUN_DEL,notEmpty" envDefault:"30m"`
}

func InitConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalln(err)
	}
	return cfg
}
