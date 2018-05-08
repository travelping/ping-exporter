package main

import (
	"flag"
	"os"
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/digineo/go-ping"

     mon "github.com/digineo/go-ping/monitor"
    "github.com/prometheus/common/log"
	"net"
	"time"
	"strings"
)

const version string = "0.1.0"

var (
    showVersion = flag.Bool("version", false, "Print version information")
    addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
    metricsPath = flag.String("web.telemetry-path", "/metrics", "Path under which metrics will be exposed")
	listenAddress = flag.String("web.listen-address", ":9427", "Address on which to expose metrics and web interface")
	pingInterval  = flag.Duration("ping.interval", time.Duration(5)*time.Second, "Interval for ICMP echo requests")
	pingTimeout   = flag.Duration("ping.timeout", time.Duration(4)*time.Second, "Timeout for ICMP echo request")
	pingTarget		  = flag.String("ping.target", "9.9.9.9", "IP address of target")
	dnsRefresh    = flag.Duration("dns.refresh", time.Duration(1)*time.Minute, "Interval for refreshing DNS records and updating targets accordingly (0 if disabled)")

)

var (
	uniformDomain     = flag.Float64("uniform.domain", 0.0002, "The domain for the uniform distribution.")
	normDomain        = flag.Float64("normal.domain", 0.0002, "The domain for the normal distribution.")
	normMean          = flag.Float64("normal.mean", 0.00001, "The mean for the normal distribution.")
)

var (
    pingDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:       "ping_durations_histogram_seconds",
        Help:       "Ping Latency distribution.",
        Buckets:    prometheus.LinearBuckets(*normMean-5**normDomain, .5**normDomain, 20),
    })
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
	pinger, err := ping.New("0.0.0.0", "::")
	if err != nil {
		return nil, err
	}

	monitor := mon.New(pinger, *pingInterval, *pingTimeout)

	targetList := strings.Split(*pingTarget, ",")

	targets := make([]*target, len(targetList))
	for i, host := range targetList {
		t := &target {
			host: host,
			addresses: make([]net.IP, 0),
			delay: time.Duration(10*i) * time.Millisecond,
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
	if *dnsRefresh == 0 {
		return
	}

	for {
		select {
		case <-time.After(*dnsRefresh):
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
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(pingDurationHistogram)
	registry.MustRegister(&pingCollector{monitor: monitor})

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError})
	http.HandleFunc(*metricsPath, h.ServeHTTP)
	log.Infof("Listening for %s on %s", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func main() {
    flag.Parse()
    if *showVersion {
        printVersion()
        os.Exit(0)
    }

    m, err := startMonitor()
    if err != nil {
    	log.Errorln(err)
    	os.Exit(2)
	}
    startServer(m)
}
