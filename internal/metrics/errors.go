package metrics

import (
	"github.com/hpcloud/tail"
	"logger/internal/helper"
	"logger/internal/options"
	"strconv"
	"sync"
)

type Errors struct {
	options options.Options
	MetricsErrors
}

type MetricsErrors struct {
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
		t, _ := tail.TailFile(a.options.ErrorsPath, tail.Config{Follow: true})
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

	if helper.Contains(scannerBytes, []string{"PHP Fatal", "PHP Parse", "PHP Core"}) &&
		!helper.Contains(scannerBytes, []string{"Maximum execution time", "SQLERR:"}) {
		a.incrementFatal()
	}

	if helper.Contains(scannerBytes, []string{"Maximum execution time"}) {
		a.incrementTime()
	}

	if helper.Contains(scannerBytes, []string{"PHP Warning"}) &&
		!helper.Contains(scannerBytes, []string{"SQLERR"}) &&
		!helper.Contains(scannerBytes, []string{"Maximum execution time"}) {
		a.incrementWarning()
	}

	if helper.Contains(scannerBytes, []string{"SQLERR"}) {
		a.incrementSql()
	}

	if helper.Contains(scannerBytes, []string{"PHP Notice", "PHP Deprecated", "PHP Strict"}) {
		a.incrementPhpNotice()
	}

	if !helper.Contains(scannerBytes, []string{"PHP", "SQLERR"}) {
		a.incrementPhpNotice()
	}
}

func (a *Errors) getAllMetrics() map[string]string {
	prefix := ""
	if a.options.Debug {
		prefix = "_test"
	}

	return map[string]string{
		"apache.errors.time" + prefix:    strconv.Itoa(a.time),
		"apache.errors.fatal" + prefix:   strconv.Itoa(a.fatal),
		"apache.errors.warning" + prefix: strconv.Itoa(a.warning),
		"apache.errors.sql" + prefix:     strconv.Itoa(a.sql),
		"apache.errors.notice" + prefix:  strconv.Itoa(a.phpNotice),
		"apache.errors.nonphp" + prefix:  strconv.Itoa(a.nonPhp),
	}
}

func (a *Errors) SetParams(o options.Options) {
	a.options = o
}
