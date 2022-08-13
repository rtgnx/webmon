package webmon

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Probe interface {
	Ping() error
}

const (
	HTTP = "http"
	TCP  = "tcp"
)

type ProbeConfig struct {
	Target   string        `json:"target,omitempty"`
	Proto    string        `json:"proto,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
}

func (cfg ProbeConfig) Probe() (Probe, error) {
	switch cfg.Proto {
	case HTTP:
		return NewHTTPProbe(cfg.Target)
	case TCP:
		return NewTCPProbe(cfg.Target)
	default:
		return nil, fmt.Errorf("unsupported protocol %s", cfg.Proto)
	}
}

func Monitor(config []ProbeConfig) {
	for _, cfg := range config {
		probe, err := cfg.Probe()
		if err != nil {
			log.Print(err)
			continue
		}

		go func(probe Probe, cfg ProbeConfig) {
			metric := promauto.NewGauge(prometheus.GaugeOpts{
				Namespace: "webmon",
				Subsystem: cfg.Proto,
				Name:      "ping",
				Help:      "",
				ConstLabels: map[string]string{
					"target":   cfg.Target,
					"proto":    cfg.Proto,
					"interval": cfg.Interval.String(),
				},
			})
			for stop := false; !stop; {
				<-time.After(cfg.Interval)
				if err := probe.Ping(); err != nil {
					log.Println(err)
					metric.Add(1.0)
					continue
				}
				metric.Set(0.0)
			}
		}(probe, cfg)
	}
}
