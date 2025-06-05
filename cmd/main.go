package main

import (
	"bytes"
	"context"
	"log"
	"net/http"

	"github.com/caarlos0/env/v11"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/fideltak/cecp/pkg/prometheus"
)

var (
	server  *Server
	version string = "dev"
)

type Server struct {
	Config   *Config
	Exporter *prometheus.Exporter
}

type Config struct {
	HttpHost  string `env:"HTTP_HOST,required" envDefault:"0.0.0.0"`
	HttpPort  int    `env:"HTTP_PORT,required" envDefault:"8080"`
	TargetURL string `env:"TARGET_URL,required" envDefault:"http://localhost:8888"`
}

func (c *Config) GetHTTPHost() string {
	if c != nil {
		return c.HttpHost
	}
	return ""
}

func (c *Config) GetHTTPPort() int {
	if c != nil {
		return c.HttpPort
	}
	return 0
}

func (c *Config) GetTargetURL() string {
	if c != nil {
		return c.TargetURL
	}
	return ""
}

func NewConfig(optConfig ...*Config) (*Config, error) {
	log.Printf("check configuration for CloudEvents Proxy.")
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	if len(optConfig) != 0 {
		cfg = optConfig[0]
	}
	return cfg, nil
}

func NewServer(serverConfig *Config, promConfig *prometheus.Config) (*Server, error) {
	var err error
	var serviceConfig *Config

	if serverConfig != nil {
		serviceConfig, err = NewConfig(serverConfig)
	} else {
		serviceConfig, err = NewConfig()
	}
	if err != nil {
		return nil, err
	}

	var exporter *prometheus.Exporter
	if promConfig != nil {
		exporter, err = prometheus.NewExporter(promConfig)
	} else {
		exporter, err = prometheus.NewExporter()
	}
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:   serviceConfig,
		Exporter: exporter,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	// Run the Prometheus exporter
	go func() {
		s.Exporter.Run()
	}()
	err := startReceiver(ctx)
	if err != nil {
		return err
	}
	return nil
}

func startReceiver(ctx context.Context) error {
	httpHost := server.Config.GetHTTPHost()
	httpPort := server.Config.GetHTTPPort()
	log.Printf("starting CloudEvents receiver on %s:%v ...\n", httpHost, httpPort)
	p, err := cloudevents.NewHTTP(cloudevents.WithPort(httpPort), cloudevents.WithHost(httpHost))
	if err != nil {
		return err
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		return err
	}
	log.Printf("listening for CloudEvents on %s:%v ...\n", httpHost, httpPort)

	// Start listening and call handleEvent for every received event
	err = c.StartReceiver(ctx, proxy)
	if err != nil {
		return err
	}
	return nil
}

func proxy(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	targetURL := server.Config.GetTargetURL()
	promProxyRequest.Inc() // Increment the proxy request counter
	jsonBytes, err := event.MarshalJSON()
	if err != nil {
		log.Printf("failed to marshal CloudEvent to JSON: %v", err)
		return nil, cloudevents.NewHTTPResult(500, "marshal error")
	}
	resp, err := http.Post(targetURL,
		"application/json",
		bytes.NewReader(jsonBytes))
	if err != nil {
		log.Printf("failed to post to remote server: %v", err)
		return nil, cloudevents.NewHTTPResult(502, "forwarding error")
	}
	defer resp.Body.Close()

	log.Printf("forwarded event to %s, status: %s", targetURL, resp.Status)
	log.Printf("%s", jsonBytes)
	return nil, cloudevents.ResultACK
}

func init() {
	var err error
	server, err = NewServer(nil, nil)
	if err != nil {
		log.Fatalf("failed to create server: %v\n", err)
	}
}

func main() {
	log.Printf("CloudEvents Proxy version: %s\n", version)
	ctx := context.Background()
	err := server.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
}
