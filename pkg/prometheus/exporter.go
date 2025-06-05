package prometheus

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	Config *Config
	URL    string
}

type Config struct {
	Address string `env:"PROMETHEUS_ADDRESS,required" envDefault:"0.0.0.0"`
	Port    uint16 `env:"PROMETHEUS_PORT,required" envDefault:"9100"`
	URLPath string `env:"PROMETHEUS_URL_PATH,required" envDefault:"/metrics"`
}

func (c *Config) GetAddress() string {
	if c != nil {
		return c.Address
	}
	return ""
}

func (c *Config) GetPort() uint16 {
	if c != nil {
		return c.Port
	}
	return 0
}

func (c *Config) GetURLPath() string {
	if c != nil {
		return c.URLPath
	}
	return ""
}

func (e *Exporter) GetConfig() *Config {
	if e != nil {
		return e.Config
	}
	return nil
}

func (e *Exporter) GetURL() string {
	if e != nil {
		return e.URL
	}
	return ""
}

func NewExporter(optConfig ...*Config) (*Exporter, error) {
	slog.Info("check configuration for Promethes.")
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	if len(optConfig) != 0 {
		config = optConfig[0]
	}
	return &Exporter{
		Config: config,
		URL:    fmt.Sprintf("%s:%v", config.GetAddress(), config.GetPort()),
	}, nil
}

func (e *Exporter) Run() {
	config := e.GetConfig()
	if config == nil {
		panic("no promehetus config")
	}
	url := e.GetURL()
	slog.Info(fmt.Sprintf("Start prometheus exporter %s", url))
	http.Handle(config.GetURLPath(), promhttp.Handler())
	go func() {
		err := http.ListenAndServe(url, nil)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()
}
