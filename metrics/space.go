package metrics

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"
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

func commandExec(search string) map[string]string {
	spaceCmd := exec.Command("df", "-B", "1M")
	findCmd := exec.Command("grep", search)

	reader, writer := io.Pipe()
	var buf bytes.Buffer
	spaceCmd.Stdout = writer
	findCmd.Stdin = reader
	findCmd.Stdout = &buf

	spaceCmd.Start()
	findCmd.Start()

	spaceCmd.Wait()
	writer.Close()

	findCmd.Wait()
	reader.Close()
	io.Copy(os.Stdout, &buf)

	s := regexp.MustCompile("/\\s\\s+/").Split(buf.String(), 5)
	if len(s) == 5 {
		return map[string]string{
			"total":        s[1],
			"used":         s[2],
			"left":         s[3],
			"percent_used": s[3],
		}
	}
	return map[string]string{}
}

func (s *Space) handleText() map[string]string {
	//res := commandExec("data-front")
	res := commandExec("loop51")
	//res1 := commandExec("storage")
	res1 := commandExec("loop52")
	return map[string]string{
		"disk.root.total":        res["total"],
		"disk.root.used":         res["used"],
		"disk.root.left":         res["left"],
		"disk.root.percent_used": res["percent_used"],

		"disk.storage.total":        res1["total"],
		"disk.storage.used":         res1["used"],
		"disk.storage.left":         res1["left"],
		"disk.storage.percent_used": res1["percent_used"],
	}
}

func (s *Space) Get(done <-chan struct{}, chanText chan map[string]string, wg *sync.WaitGroup) {
	wg.Add(1)
	func() {
		ticker := time.NewTicker(time.Second)
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				chanText <- s.handleText()
			}
		}
	}()
	wg.Wait()
}
