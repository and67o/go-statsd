package metrics

import (
	"bytes"
	"logger/internal/options"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Space struct {
	options options.Options

	Root
	Storage
}

type Root struct {
	Data
}
type Storage struct {
	Data
}

type Data struct {
	total       string
	used        string
	left        string
	percentUsed string
}

func (d *Data) SetFront(data map[string]string) {
	usedInt, _ := strconv.Atoi(data["used"])
	totalInt, _ := strconv.Atoi(data["total"])
	percentUsed := 0
	if totalInt > 0 {
		percentUsed = (usedInt / totalInt) * 100
	}

	d.used = data["used"]
	d.total = data["total"]
	d.left = string(rune(totalInt - usedInt))
	d.percentUsed = string(rune(percentUsed))
}

func commandExec(search string) map[string]string {
	spaceCmd := exec.Command("df", "-B", "1M")
	grepCmd := exec.Command("grep", search)
	var buf bytes.Buffer

	grepCmd.Stdin, _ = spaceCmd.StdoutPipe()
	grepCmd.Stdout = os.Stdout
	grepCmd.Stdout = &buf

	_ = grepCmd.Start()
	_ = spaceCmd.Run()
	_ = grepCmd.Wait()

	resGrep := strings.Fields(buf.String())
	if len(resGrep) == 6 {
		return map[string]string{
			"total":        resGrep[1],
			"used":         resGrep[2],
			"left":         resGrep[3],
			"percent_used": resGrep[4],
		}
	}
	return map[string]string{}
}

func (s *Space) handle() {
	dataFront := commandExec("data-front") //"data-front"
	dataStorage := commandExec("storage")  //"storage"

	s.Root.SetFront(dataFront)
	s.Storage.SetFront(dataStorage)
}

func (s *Space) getAllMetrics() map[string]string {
	prefix := ""
	if s.options.Debug {
		prefix = "_test"
	}

	return map[string]string{
		"disk.root.total" + prefix:        s.Root.total,
		"disk.root.used" + prefix:         s.Root.used,
		"disk.root.left" + prefix:         s.Root.left,
		"disk.root.percent_used" + prefix: s.Root.percentUsed,

		"disk.storage.total" + prefix:        s.Storage.total,
		"disk.storage.used" + prefix:         s.Storage.used,
		"disk.storage.left" + prefix:         s.Storage.left,
		"disk.storage.percent_used" + prefix: s.Storage.percentUsed,
	}
}

func (s *Space) Get(done <-chan struct{}, chanText chan map[string]string, wg *sync.WaitGroup) {
	wg.Add(1)
	func() {
		ticker := time.NewTicker(time.Second)
		defer wg.Done()
		for {
			s.handle()
			select {
			case <-done:
				return
			case <-ticker.C:
				chanText <- s.getAllMetrics()
			}
		}
	}()
	wg.Wait()
}

func (s *Space) SetParams(o options.Options) {
	s.options = o
}
