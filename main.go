package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	ping "github.com/digineo/go-ping"
	mon "github.com/digineo/go-ping/monitor"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

var (
	showVersion = flag.Bool("version", false, "Print version information")
)

const version = "0.4.0"

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
	fmt.Println("ping-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Tobias Famulla")
}

func startMonitor(config PingConfig, dnsRefresh time.Duration) (*mon.Monitor, error) {
	pinger, err := ping.New(config.SourceV4, config.SourceV6)
	if err != nil {
		return nil, err
	}

	monitor := mon.New(pinger, config.PingInterval, config.PingTimeout)
	targets := []*pingTarget(nil)

	for i, host := range config.PingTargets {
		t := &pingTarget{
			host:      host,
			addresses: make([]net.IP, 0),
			delay:     time.Duration(10*i) * time.Millisecond,
			sourceV4:  config.SourceV4,
			sourceV6:  config.SourceV6,
		}
		targets = append(targets, t)

		err := t.addOrUpdateMonitor(monitor)
		if err != nil {
			log.Errorln(err)
		}
	}

	go startDNSAutoRefresh(dnsRefresh, targets, monitor)

	return monitor, nil

}

func startDNSAutoRefresh(dnsRefresh time.Duration, targets []*pingTarget, monitor *mon.Monitor) {
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

func refreshDNS(targets []*pingTarget, monitor *mon.Monitor) {
	for _, t := range targets {
		log.Infof("refreshing DNS")

		go func(ta *pingTarget) {
			err := ta.addOrUpdateMonitor(monitor)
			if err != nil {
				log.Errorf("could refresh dns: %v", err)
			}
		}(t)
	}
}

func startServer(monitor []*mon.Monitor, metricsPath string, listenAddress string) {
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
	registry.MustRegister(&pingCollector{monitors: monitor})

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

	config := Configuration{}
	currentViper := viper.New()

	initViper(currentViper)
	readInConfig(currentViper)
	config.updateConfig(currentViper)

	configSet, missingParameters := isMandatoryConfigSet(currentViper)
	if !configSet {
		log.Errorln("configuration parameters", missingParameters, "missing")
		os.Exit(3)
	}

	var monitors []*mon.Monitor

	if config.hasPingMultiConfig {
		for _, c := range config.pingConfigurations {
			m, err := startMonitor(c, config.dnsRefresh)
			if err != nil {
				log.Errorln(err)
				os.Exit(2)
			}
			monitors = append(monitors, m)
		}

	} else {
		target := PingConfig{config.pingSourceV4, config.pingSourceV6, config.pingTarget, config.pingInterval, config.pingTimeout}
		m, err := startMonitor(target, config.dnsRefresh)
		if err != nil {
			log.Errorln(err)
			os.Exit(2)
		}
		monitors = append(monitors, m)
	}

	startServer(monitors, config.metricsPath, config.listenAddress)
}
