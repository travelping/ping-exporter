package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	mon "github.com/digineo/go-ping/monitor"
	"github.com/prometheus/common/log"
)

type pingTarget struct {
	host         string
	addresses    []net.IP
	delay        time.Duration
	mutex        sync.Mutex
	sourceV4     string
	sourceV6     string
	sourceLabels []keyValuePair
	targetLabels []keyValuePair
}

type keyValuePair struct {
	key   string
	value string
}

func (t *pingTarget) addOrUpdateMonitor(monitor *mon.Monitor) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	addrs, err := net.LookupIP(t.host)
	if err != nil {
		return fmt.Errorf("error resolving target: %v", err)
	}

	for _, addr := range addrs {
		err := t.addIfNew(addr, monitor)
		if err != nil {
			return err
		}
	}

	t.cleanUp(addrs, monitor)
	t.addresses = addrs

	return nil
}

func (t *pingTarget) addIfNew(addr net.IP, monitor *mon.Monitor) error {
	if isIPInSlice(addr, t.addresses) {
		return nil
	}

	return t.add(addr, monitor)
}

func (t *pingTarget) cleanUp(new []net.IP, monitor *mon.Monitor) {
	for _, o := range t.addresses {
		if !isIPInSlice(o, new) {
			name := t.nameForIP(o)
			log.Infof("removing target for host %s (%v)", t.host, o)
			monitor.RemoveTarget(name)
		}
	}
}

func (t *pingTarget) add(addr net.IP, monitor *mon.Monitor) error {
	name := t.nameForIP(addr)
	log.Infof("adding target for host %s (%v)", t.host, addr)
	return monitor.AddTargetDelayed(name, net.IPAddr{IP: addr, Zone: ""}, t.delay)
}

func (t *pingTarget) nameForIP(addr net.IP) string {
	name := fmt.Sprintf("%s %s ", t.host, addr)

	if addr.To4() == nil {
		name += fmt.Sprintf("6 %s", t.sourceV6)
	} else {
		name += fmt.Sprintf("4 %s", t.sourceV4)
	}

	return name
}

func isIPInSlice(ip net.IP, slice []net.IP) bool {
	for _, x := range slice {
		if x.Equal(ip) {
			return true
		}
	}
	return false
}
