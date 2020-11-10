package metrics

import "fmt"

type Errors struct {}

func (a *Errors) Set(name string, value interface{}) {
	fmt.Println(333)
}

func (a *Errors) Get(name string) map[string]int {
	return map[string]int{
		"pache.requests.total1":    5,
		"apache.requests.api1":    6,
		"apache.requests.storage1": 7,
		"apache.requests.site1":    8,
	}
}
