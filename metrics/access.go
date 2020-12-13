package metrics

import (
	"github.com/hpcloud/tail"
	"logger/bytesUtils"
	"strconv"
	"sync"
)

const file = "/var/log/apache2/dev.eljur.access.log"

type Access struct {
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
		t, _ := tail.TailFile(file, tail.Config{Follow: true})
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
	return map[string]string{
		"apache.requests.total":   strconv.Itoa(a.total),
		"apache.requests.api":     strconv.Itoa(a.api),
		"apache.requests.storage": strconv.Itoa(a.storage),
		"apache.requests.site":    strconv.Itoa(a.other),
	}
}

func (a *Access) handleText(scannerBytes []byte) {
	a.incrementTotal()

	if bytesUtils.Contains(scannerBytes, []string{"/api"}) {
		a.incrementApi()
	}

	if bytesUtils.Contains(scannerBytes, []string{"/storage"}) {
		a.incrementStorage()
	}

	if !bytesUtils.Contains(scannerBytes, []string{"/api", "/storage"}) {
		a.incrementOther()
	}
}
