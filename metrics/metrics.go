package metrics

import "sync"

type Metrics interface {
	Get(done <-chan struct{}, chanText chan map[string]string, wg *sync.WaitGroup)
}
