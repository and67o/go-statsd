package metrics

import "fmt"

type Space struct {}

func (a *Space) Set(name string, value interface{}) {
	fmt.Println(444)
}

func (a *Space) Get(name string) (interface{}, error) {
	return 5, nil
}

