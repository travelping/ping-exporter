package main

import (
	"strings"
	"sync"

	mon "github.com/digineo/go-ping/monitor"
	"github.com/imdario/mergo"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "ping_"

var (
	labelNames = []string{"target", "ip", "ip_version", "source_ip"}
	mutex      = &sync.Mutex{}
)

func rttDesc(addLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prefix+"rtt_ms", "Round trip time in millis (deprecated)", append(append(labelNames, addLabels...), "type"), nil)
}

func bestDesc(addLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prefix+"rtt_best_ms", "Best round trip time in millis", append(labelNames, addLabels...), nil)
}

func worstDesc(addLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prefix+"rtt_worst_ms", "Worst round trip time in millis", append(labelNames, addLabels...), nil)
}

func meanDesc(addLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prefix+"rtt_mean_ms", "Mean round trip time in millis", append(labelNames, addLabels...), nil)
}

func stddevDesc(addLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prefix+"rtt_std_deviation_ms", "Standard deviation in millis", append(labelNames, addLabels...), nil)
}

func lossDesc(addLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(prefix+"loss_percent", "Packet loss in percent", append(labelNames, addLabels...), nil)
}

type pingCollector struct {
	pingMonitors []*pingMonitor
	metrics      map[string]*mon.Metrics
}

func findAdditionalLabels(target string, source string, ip_version string, p *pingCollector) ([]string, []string, error) {
	var targetLabels []keyValuePair
	var sourceLabels []keyValuePair
	for _, pm := range p.pingMonitors {
		for _, t := range *pm.targets {
			if ip_version == "4" {
				if t.host == target && t.sourceV4 == source {
					targetLabels = t.targetLabels
					sourceLabels = t.sourceLabels
				}
			} else {
				if t.host == target && t.sourceV6 == source {
					targetLabels = t.targetLabels
					sourceLabels = t.sourceLabels
				}
			}
		}
	}
	keys := []string{}
	values := []string{}

	for _, p := range sourceLabels {
		keys = append(keys, p.key)
		values = append(values, p.value)
	}
	for _, p := range targetLabels {
		keys = append(keys, p.key)
		values = append(values, p.value)
	}
	return keys, values, nil
}

func (p *pingCollector) Describe(ch chan<- *prometheus.Desc) {
	addLabels := []string{}
	ch <- rttDesc(addLabels)
	ch <- lossDesc(addLabels)
	ch <- bestDesc(addLabels)
	ch <- worstDesc(addLabels)
	ch <- meanDesc(addLabels)
	ch <- stddevDesc(addLabels)
}

func (p *pingCollector) Collect(ch chan<- prometheus.Metric) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, pm := range p.pingMonitors {
		metrics := pm.monitor.ExportAndClear()

		if len(metrics) > 0 {
			err := mergo.Merge(&p.metrics, metrics, mergo.WithOverride)
			if err != nil {
				panic(err)
			}
		}

	}

	if p.metrics == nil || len(p.metrics) == 0 {
		return
	}

	for target, metrics := range p.metrics {
		t := strings.Split(target, " ")
		l := []string{t[0], t[1], t[2], t[3]}

		labels, labelValues, _ := findAdditionalLabels(t[0], t[3], t[2], p)

		l = append(l, labelValues...)

		ch <- prometheus.MustNewConstMetric(rttDesc(labels), prometheus.GaugeValue, float64(metrics.Best), append(l, "best")...)
		ch <- prometheus.MustNewConstMetric(bestDesc(labels), prometheus.GaugeValue, float64(metrics.Best), l...)

		ch <- prometheus.MustNewConstMetric(rttDesc(labels), prometheus.GaugeValue, float64(metrics.Worst), append(l, "worst")...)
		ch <- prometheus.MustNewConstMetric(worstDesc(labels), prometheus.GaugeValue, float64(metrics.Worst), l...)

		ch <- prometheus.MustNewConstMetric(rttDesc(labels), prometheus.GaugeValue, float64(metrics.Mean), append(l, "mean")...)
		ch <- prometheus.MustNewConstMetric(meanDesc(labels), prometheus.GaugeValue, float64(metrics.Mean), l...)

		ch <- prometheus.MustNewConstMetric(rttDesc(labels), prometheus.GaugeValue, float64(metrics.StdDev), append(l, "std_dev")...)
		ch <- prometheus.MustNewConstMetric(stddevDesc(labels), prometheus.GaugeValue, float64(metrics.StdDev), l...)

		loss := float64(metrics.PacketsLost) / float64(metrics.PacketsSent)
		ch <- prometheus.MustNewConstMetric(lossDesc(labels), prometheus.GaugeValue, loss, l...)
	}
}
