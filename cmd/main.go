package main

import (
	"flag"
	"logger/internal/options"
)

var (
	host       string
	port       string
	accessPath string
	errorPath  string
	nameServer string
	prefix     string
	debug      bool
)

func init() {
	//flag.StringVar(&host, "host", "127.0.0.1", "host")
	//flag.StringVar(&port, "port", "8125", "host")
	//flag.StringVar(&accessPath, "access-path", "/var/log/apache2/dev.eljur.access.log", "path to access log")
	//flag.StringVar(&errorPath, "error-path", "/var/log/apache2/dev.eljur.error.log", "path to error log")

	flag.StringVar(&host, "host", "is-eljur-ok.ru", "host")
	flag.StringVar(&port, "port", "8125", "host")
	flag.StringVar(&accessPath, "access-path", "/home/webadmin/eljur.ru/logs/www.access.log", "path to access log")
	flag.StringVar(&errorPath, "error-path", "/home/webadmin/eljur.ru/logs/www.error.log", "path to error log")

	flag.StringVar(&nameServer, "name-server", "front-1", "server name")
	flag.StringVar(&prefix, "prefix", "prod", "prefix")
	flag.BoolVar(&debug, "debug", true, "add prefix to metrics")
}

func main() {
	flag.Parse()

	a := App{}

	a.options = initOptions()

	a.Initialize()

	a.Run()
}

func initOptions() options.Options {
	var o options.Options

	o.Host = host
	o.Port = port
	o.AccessPath = accessPath
	o.ErrorsPath = errorPath
	o.Name = nameServer
	o.Prefix = prefix
	o.Debug = debug

	return o
}
