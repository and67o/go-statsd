package apacheRequests

import (
	"logger/requestUtils"
)

const Total = "apache.requests.total"
const Api = "apache.requests.api"
const Storage = "apache.requests.storage"
const Site = "apache.requests.site"

type totalCommandsMetrics struct {
	total, api, storage, other int
	result                     map[string]int
}

func (r *totalCommandsMetrics) addMetrics(key string, value int) {
	r.result[key] = value
}

type requestInterface interface {
	Get() map[string]int
}

func GetTotalApi() (int, error) {
	res, err := requestUtils.SubStringCommand("api")
	if err != nil {
		return 0, err
	}
	return res, nil
}

func Get() map[string]int {
	var res make(map[string]int)

}
