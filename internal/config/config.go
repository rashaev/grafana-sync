package config

import (
	"log"
	"net/url"
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	GrafanaURL           url.URL       `env:"GRAFANA_SYNC_GRAFANA_URL,required"`
	GrafanaUser          string        `env:"GRAFANA_SYNC_GRAFANA_USER,required"`
	GrafanaPassword      string        `env:"GRAFANA_SYNC_GRAFANA_PASSWORD,required"`
	AvanpostURL          url.URL       `env:"GRAFANA_SYNC_AVANPOST_URL,required"`
	AvanpostUser         string        `env:"GRAFANA_SYNC_AVANPOST_USER,required"`
	AvanpostPassword     string        `env:"GRAFANA_SYNC_AVANPOST_PASSWORD,required"`
	AvanpostClientID     string        `env:"GRAFANA_SYNC_AVANPOST_CLIENT_ID,required"`
	AvanpostClientSecret string        `env:"GRAFANA_SYNC_AVANPOST_CLIENT_SECRET,required"`
	GroupRegexRO         string        `env:"GRAFANA_SYNC_GROUP_REGEX_RO,required"`
	GroupRegexRW         string        `env:"GRAFANA_SYNC_GROUP_REGEX_RW,required"`
	TimeRunSync          time.Duration `env:"GRAFANA_SYNC_TIME_RUN_SYNC,required"`
}

func InitConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalln(err)
	}
	return cfg
}
