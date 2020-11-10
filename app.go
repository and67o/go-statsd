package main

import (
	"logger/metrics"
	"logger/statsD"
)

type App struct {
	StatsD statsD.Operations
}

func (a *App) Initialize() {
	a.StatsD = statsD.New()
}

func (a *App) Run() {
	m := getAllMetrics()

	c := a.StatsD
	defer c.Close()

	done := make(chan interface{})
	defer close(done)

	intStream := worker(done, m)
	for v := range intStream {
		for key, value := range v {
			c.Gauge(key, value)
		}
	}
}

func getAllMetrics() []metrics.Metrics {
	return []metrics.Metrics{
		&metrics.Access{},
		&metrics.Errors{},
		//&metrics.Space{},
	}
}

func worker(done <-chan interface{}, m []metrics.Metrics) <-chan map[string]int {
	intStream := make(chan map[string]int)
	go func() {
		defer close(intStream)
		for _, metric := range m {
			select {
			case <-done:
				return
			case intStream <- metric.Get("oleg"):
			}
		}
	}()
	return intStream
}
