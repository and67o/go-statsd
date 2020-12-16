package main

import (
	"fmt"
	"logger/metrics"
	"logger/statsD"
	"strconv"
	"sync"
)

type App struct {
	StatsD statsD.Operations
}

func (a *App) Initialize() {
	a.StatsD = statsD.New()
}

func (a *App) Run(stopCh chan struct{}) {
	c := a.StatsD
	defer c.Close()

	m := getAllMetrics()

	metricsStr := worker(stopCh, m)

	for metricStr := range metricsStr {
		for nameMetric, text := range metricStr {
			v, _ := strconv.Atoi(text)
			c.Gauge(nameMetric, v)
			fmt.Println(nameMetric, text)
		}
	}
}

func worker(stopCh <-chan struct{}, metrics []metrics.Metrics) <-chan map[string]string {
	textChan := make(chan map[string]string)
	wg := sync.WaitGroup{}
	for _, metric := range metrics {
		go metric.Get(stopCh, textChan, &wg)
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
