package metrics

import (
	"github.com/hpcloud/tail"
	"logger/internal/helper"
	"logger/internal/options"
	"strconv"
	"sync"
)

type Access struct {
	options options.Options

	total   int
	api     int
	storage int
	other   int
}

func (a *Access) incrementTotal() {
	a.total++
}

func (a *Access) incrementApi() {
	a.api++
}

func (a *Access) incrementStorage() {
	a.storage++
}

func (a *Access) incrementOther() {
	a.other++
}

func (a *Access) Get(done <-chan struct{}, chanText chan map[string]string, wg *sync.WaitGroup) {
	wg.Add(1)
	func() {
		t, _ := tail.TailFile(a.options.AccessPath, tail.Config{Follow: true})
		defer wg.Done()
		for line := range t.Lines {
			a.handleText([]byte(line.Text))
			select {
			case <-done:
				return
			case chanText <- a.getAllMetrics():
			}
		}
	}()
	wg.Wait()
}

func (a *Access) getAllMetrics() map[string]string {
	prefix := ""
	if a.options.Debug {
		prefix = "_test"
	}

	return map[string]string{
		"apache.requests.total" + prefix:   strconv.Itoa(a.total),
		"apache.requests.api" + prefix:     strconv.Itoa(a.api),
		"apache.requests.storage" + prefix: strconv.Itoa(a.storage),
		"apache.requests.site" + prefix:    strconv.Itoa(a.other),
	}
}

func (a *Access) handleText(scannerBytes []byte) {
	a.incrementTotal()

	if helper.Contains(scannerBytes, []string{"/api"}) {
		a.incrementApi()
	}

	if helper.Contains(scannerBytes, []string{"/storage"}) {
		a.incrementStorage()
	}

	if !helper.Contains(scannerBytes, []string{"/api", "/storage"}) {
		a.incrementOther()
	}
}

func (a *Access) SetParams(o options.Options) {
	a.options = o
}
