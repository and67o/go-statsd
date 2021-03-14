package main

import (
	"logger/internal/metrics"
	"logger/internal/options"
	"logger/internal/statsD"
	"strconv"
	"sync"
)

type App struct {
	StatsD  statsD.Operations
	options options.Options
}

func (a *App) Initialize() {
	a.StatsD = statsD.New(a.options)
}

func (a *App) Run() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	c := a.StatsD
	defer c.Close()

	metricsStr := a.worker(stopCh)

	for metricStr := range metricsStr {
		for nameMetric, text := range metricStr {
			v, _ := strconv.Atoi(text)
			c.Gauge(nameMetric, v)
		}
	}
}

func (a *App) worker(stopCh <-chan struct{}) <-chan map[string]string {
	m := getAllMetrics()

	textChan := make(chan map[string]string)

	var wg sync.WaitGroup

	for _, metric := range m {
		metric.SetParams(a.options)

		go metric.Get(
			stopCh,
			textChan,
			&wg,
		)
	}

	return textChan
}

func getAllMetrics() []metrics.Metrics {
	return []metrics.Metrics{
		&metrics.Access{},
		&metrics.Errors{},
		&metrics.Space{},
	}
}
