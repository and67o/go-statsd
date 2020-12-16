package metrics

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Space struct {
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
	d.left = data["left"]
	d.total = data["total"]
	d.used = data["used"]
	d.percentUsed = data["percent_used"]
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

	resGrep:= strings.Fields(buf.String())
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
	dataFront := commandExec("loop51")   //"data-front"
	dataStorage := commandExec("loop52") //"storage"

	s.Root.SetFront(dataFront)
	s.Storage.SetFront(dataStorage)
}

func (s *Space) getAllMetrics() map[string]string {
	return map[string]string{
		"disk.root.total":        s.Root.total,
		"disk.root.used":         s.Root.used,
		"disk.root.left":         s.Root.left,
		"disk.root.percent_used": s.Root.percentUsed,

		"disk.storage.total":        s.Storage.total,
		"disk.storage.used":         s.Storage.used,
		"disk.storage.left":         s.Storage.left,
		"disk.storage.percent_used": s.Storage.percentUsed,
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
