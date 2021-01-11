package main

import (
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

func (a *App) Run() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	c := a.StatsD
	defer c.Close()

	metricsStr := worker(stopCh)

	for metricStr := range metricsStr {
		for nameMetric, text := range metricStr {
			v, _ := strconv.Atoi(text)
			c.Gauge(nameMetric, v)
			//fmt.Println(nameMetric, text)
		}
	}
}

func worker(stopCh <-chan struct{}) <-chan map[string]string {
	m := getAllMetrics()
	textChan := make(chan map[string]string)
	wg := sync.WaitGroup{}
	for _, metric := range m {
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