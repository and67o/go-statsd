package main

import (
	"fmt"
	"logger/metrics"
	"logger/statsD"
)

type App struct {
	StatsD *statsD.Manager
}

func (a *App) Initialize() {
	a.StatsD = statsD.New()
}

func (a *App) Run() {
	a.StatsD.Gauge("oleg", 123)
	m := [2]metrics.Metrics{
		&metrics.Access{},
		&metrics.Errors{},
		//&metrics.Space{},
	}

	generator := func(done <-chan interface{}, m [2]metrics.Metrics) <-chan map[string]int {
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

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, m)
	for v := range intStream {
		for key, value :=range v {

			a.StatsD.Gauge(key, value)
			fmt.Println(key, value)
		}
	}
}
