package main

import (
	"flag"
	"fmt"
	"github.com/digineo/go-ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"

	mon "github.com/digineo/go-ping/monitor"
	"github.com/prometheus/common/log"
	"net"
	"time"
)

const version string = "0.1.0"

var (
	showVersion = flag.Bool("version", false, "Print version information")
)

var (
	uniformDomain = flag.Float64("uniform.domain", 0.0002, "The domain for the uniform distribution.")
	normDomain    = flag.Float64("normal.domain", 0.0002, "The domain for the normal distribution.")
	normMean      = flag.Float64("normal.mean", 0.00001, "The mean for the normal distribution.")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage:", os.Args[0], "-config.path=$my-config-file [options]")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func printVersion() {
	fmt.Println("cgw-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Tobias Famulla")
	fmt.Println("Metric exporter for Travelping CGW")
}

func startMonitor() (*mon.Monitor, error) {
	pinger, err := ping.New(pingSourceV4, pingSourceV6)
	if err != nil {
		return nil, err
	}

	monitor := mon.New(pinger, pingInterval, pingTimeout)

	targets := make([]*target, len(pingTarget))
	for i, host := range pingTarget {
		t := &target{
			host:      host,
			addresses: make([]net.IP, 0),
			delay:     time.Duration(10*i) * time.Millisecond,
		}
		targets[i] = t

		err := t.addOrUpdateMonitor(monitor)
		if err != nil {
			log.Errorln(err)
		}
	}

	go startDNSAutoRefresh(targets, monitor)

	return monitor, nil

}

func startDNSAutoRefresh(targets []*target, monitor *mon.Monitor) {
	if dnsRefresh == 0 {
		return
	}

	for {
		select {
		case <-time.After(dnsRefresh):
			refreshDNS(targets, monitor)
		}
	}
}

func refreshDNS(targets []*target, monitor *mon.Monitor) {
	for _, t := range targets {
		log.Infof("refreshing DNS")

		go func(ta *target) {
			err := ta.addOrUpdateMonitor(monitor)
			if err != nil {
				log.Errorf("could refresh dns: %v", err)
			}
		}(t)
	}
}

func startServer(monitor *mon.Monitor) {
	log.Infof("starting cgw exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>CGW Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>CGW Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(&pingCollector{monitor: monitor})

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError})
	http.HandleFunc(metricsPath, h.ServeHTTP)
	log.Infof("Listening for %s on %s", metricsPath, listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	initConfig()
	readInConfig()
	updateConfig()

	configSet, missingParameters := isMandatoryConfigSet()
	if !configSet {
		log.Errorln("configuration parameters", missingParameters, "missing")
		os.Exit(3)
	}

	m, err := startMonitor()
	if err != nil {
		log.Errorln(err)
		os.Exit(2)
	}
	startServer(m)
}
