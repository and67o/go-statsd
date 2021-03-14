package statsD

import (
	"golang.org/x/exp/errors/fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"log"
	"logger/internal/options"
	"net"
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

func New(options options.Options) *Manager {
	client, err := setClient(options)
	if err != nil {
		panic(err)
	}

	log.Print("StatsD connect")
	return &Manager{client}
}

func setClient(options options.Options) (*statsd.Client, error) {
	return statsd.New(
		statsd.Address(getAddr(options)),
		statsd.Prefix(getPrefix(options)),
	)
}

func getPrefix(o options.Options) string {
	return fmt.Sprintf("%s.%s", o.Prefix, o.Name)
}

func getAddr(options options.Options) string {
	return net.JoinHostPort(options.Host, options.Port)
}
