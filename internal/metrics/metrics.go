package metrics

import (
	"logger/internal/options"
	"sync"
)

type Metrics interface {
	SetParams(o options.Options)
	Get(done <-chan struct{}, chanText chan map[string]string, wg *sync.WaitGroup)
	getAllMetrics() map[string]string
}
