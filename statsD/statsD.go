package statsD

import (
	"golang.org/x/exp/errors"
	"gopkg.in/alexcesaro/statsd.v2"
	"log"
	"logger/config"
	"strconv"
)

//type Operations interface {
//	Gauge(bucket string, value int)
//}

type Manager struct {
	statsD *statsd.Client
}

func (M *Manager) Gauge(bucket string, value int) {
	M.statsD.Gauge(bucket, value)
}

//var StatsD Operations

func New() *Manager {
	address,err := getAddr()
	if err != nil {
		panic(err)
	}

	client, err := statsd.New(statsd.Address(address))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	//client.Gauge("num_goroutine", 321)

	log.Print("StatsD connect")
	return &Manager{client}
	//StatsD = &Manager{client}
}


func getAddr() (string, error) {
	conf := config.New()

	host := conf.Stats.Host
	if host == "" {
		return "", errors.New("no host")
	}

	port := conf.Stats.Port
	if port == 0 {
		return "", errors.New("no port")
	}
	return host + ":" + strconv.Itoa(port), nil
}