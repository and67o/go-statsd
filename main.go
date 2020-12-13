package main

import (
	"time"
)

func main() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			a := App{}
			a.Initialize()

			a.Run(stopCh)
		}

	}
}