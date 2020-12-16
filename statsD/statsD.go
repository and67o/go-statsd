package statsD

import (
	"golang.org/x/exp/errors"
	"gopkg.in/alexcesaro/statsd.v2"
	"log"
	"logger/config"
	"strconv"
)

type Operations interface {
	Gauge(bucket string, value int)
	Close()
}

type Manager struct {
	statsD *statsd.Client
}
func (M *Manager) Close() {
	M.statsD.Close()
}

func (M *Manager) Gauge(bucket string, value int) {
	M.statsD.Gauge(bucket, value)
}

func New() *Manager {
	address,err := getAddr()
	if err != nil {
		panic(err)
	}

	client, err := statsd.New(statsd.Address(address))
	if err != nil {
		panic(err)
	}

	log.Print("StatsD connect")
	return &Manager{client}
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