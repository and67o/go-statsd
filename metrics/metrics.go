package metrics

type Metrics interface {
	//Get(name string) (interface{}, error)
	Get(name string) map[string]int
	Set(name string, value interface{})
}
