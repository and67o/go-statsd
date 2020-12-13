package metrics

import (
	"github.com/hpcloud/tail"
	"logger/bytesUtils"
	"strconv"
	"sync"
)

const fileErrors = "/var/log/apache2/dev.eljur.error.log"

type Errors struct {
	fatal     int
	time      int
	warning   int
	sql       int
	phpNotice int
	nonPhp    int
}

func (a *Errors) incrementNonPhp() {
	a.nonPhp++
}

func (a *Errors) incrementPhpNotice() {
	a.phpNotice++
}

func (a *Errors) incrementSql() {
	a.sql++
}

func (a *Errors) incrementWarning() {
	a.warning++
}
func (a *Errors) incrementTime() {
	a.time++
}

func (a *Errors) incrementFatal() {
	a.fatal++
}

func (a *Errors) Get(done <-chan struct{}, chanText chan map[string]string, wg *sync.WaitGroup) {
	wg.Add(1)
	func() {
		t, _ := tail.TailFile(fileErrors, tail.Config{Follow: true})
		defer wg.Done()
		for v := range t.Lines {
			a.handleText([]byte(v.Text))
			select {
			case <-done:
				return
			case chanText <- a.getAllMetrics():
			}
		}
	}()
	wg.Wait()
}

func (a *Errors) handleText(scannerBytes []byte) {

	if bytesUtils.Contains(scannerBytes, []string{"PHP Fatal", "PHP Parse", "PHP Core"}) &&
		!bytesUtils.Contains(scannerBytes, []string{"Maximum execution time", "SQLERR:"}) {
		a.incrementFatal()
	}

	if bytesUtils.Contains(scannerBytes, []string{"Maximum execution time"}) {
		a.incrementTime()
	}

	if bytesUtils.Contains(scannerBytes, []string{"PHP Warning"}) {
		a.incrementWarning()
	}

	if bytesUtils.Contains(scannerBytes, []string{"SQLERR:"}) {
		a.incrementSql()
	}

	if bytesUtils.Contains(scannerBytes, []string{"PHP Notice", "PHP Deprecated", "PHP Strict"}) {
		a.incrementPhpNotice()
	}

	if !bytesUtils.Contains(scannerBytes, []string{"PHP", "SQLERR"}) {
		a.incrementPhpNotice()
	}
}

func (a *Errors) getAllMetrics() map[string]string {
	return map[string]string{
		"apache.errors.time":    strconv.Itoa(a.time),
		"apache.errors.fatal":   strconv.Itoa(a.fatal),
		"apache.errors.warning": strconv.Itoa(a.warning),
		"apache.errors.sql":     strconv.Itoa(a.sql),
		"apache.errors.notice":  strconv.Itoa(a.phpNotice),
		"apache.errors.nonphp":  strconv.Itoa(a.nonPhp),
	}
}
