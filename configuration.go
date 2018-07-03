package main

import (
	"github.com/spf13/viper"
	//flag "github.com/spf13/pflag"
	"fmt"
	"github.com/prometheus/common/log"
	"strings"
	"time"
)

//var (
//metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which metrics will be exposed")
//listenAddress = flag.String("web.listen-address", ":9427", "Address on which to expose metrics and web interface")
//pingInterval  = flag.Duration("ping.interval", time.Duration(5)*time.Second, "Interval for ICMP echo requests")
//pingTimeout   = flag.Duration("ping.timeout", time.Duration(4)*time.Second, "Timeout for ICMP echo request")
//pingTarget    = flag.String("ping.target", "9.9.9.9", "IP address of target")
//pingSourceV4  = flag.String("ping.source.ipv4", "0.0.0.0", "Source Address of ICMP echo requests")
//pingSourceV6  = flag.String("ping.source.ipv6", "::", "Source Address of ICMP echo requests")
//dnsRefresh    = flag.Duration("dns.refresh", time.Duration(1)*time.Minute, "Interval for refreshing DNS records and updating targets accordingly (0 if disabled)")
//)
type Configuration struct {
	metricsPath        string
	listenAddress      string
	pingInterval       time.Duration
	pingTimeout        time.Duration
	pingTarget         []string
	pingSourceV4       string
	pingSourceV6       string
	hasPingMultiConfig bool
	pingConfigurations []PingConfig
	dnsRefresh         time.Duration
}

var (
	mandatoryConfig = [...]string{"web.listen-address", "ping.interval", "ping.timeout", "ping.source.ipv4", "ping.source.ipv6", "dns.refresh", "web.telemetry-path"}
	optionalConfig  = [...]string{"ping.target"}
)

type PingConfig struct {
	SourceV4     string        // source address of ICMP requests
	SourceV6     string        // source address of ICMP requests
	PingTargets  []string      // target addresses of ICMP requests
	PingInterval time.Duration // interval between ICMP requests
	PingTimeout  time.Duration // timeout of ICMP requests
}

func initViper(v *viper.Viper) {
	setDefaults(v)
	v.SetConfigName("ping-exporter")
	v.AddConfigPath("/etc/ping-exporter/")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("PINGEXPORTER")
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	bindEnvVariables(v)
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("web.listen-address", ":9427")
	v.SetDefault("web.telemetry-path", "/metrics")
	v.SetDefault("ping.interval", "5s")
	v.SetDefault("ping.timeout", "4s")
	v.SetDefault("ping.source.ipv4", "0.0.0.0")
	v.SetDefault("ping.source.ipv6", "::")
	v.SetDefault("dns.refresh", "1m")
}

func (conf *Configuration) updateConfig(v *viper.Viper) error {
	//readInConfig()
	conf.listenAddress = v.GetString("web.listen-address")
	conf.metricsPath = v.GetString("web.telemetry-path")
	conf.pingInterval = v.GetDuration("ping.interval")
	conf.pingTimeout = v.GetDuration("ping.timeout")
	conf.pingTarget = v.GetStringSlice("ping.target")
	conf.pingSourceV4 = v.GetString("ping.source.ipv4")
	conf.pingSourceV6 = v.GetString("ping.source.ipv6")
	conf.dnsRefresh = v.GetDuration("dns.refresh")
	if v.IsSet("ping.configurations") {
		conf.hasPingMultiConfig = true
		err := v.UnmarshalKey("ping.configurations", &(conf.pingConfigurations))
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
			return err
		}
	} else {
		conf.hasPingMultiConfig = false
	}
	return nil
}

func bindEnvVariables(v *viper.Viper) {
	//flag.Parse()
	//viper.BindPFlags(flag.CommandLine)

	for _, element := range mandatoryConfig {
		v.BindEnv(element)
	}
	for _, element := range optionalConfig {
		v.BindEnv(element)
	}
}

func isMandatoryConfigSet(v *viper.Viper) (bool, []string) {
	allSet := true
	var missingConfig []string
	for _, element := range mandatoryConfig {
		if !v.IsSet(element) {
			allSet = false
			missingConfig = append(missingConfig, element)
		}
	}
	if !v.IsSet("ping.target") && !v.IsSet("ping.configurations") {
		allSet = false
		missingConfig = append(missingConfig, "ping.configurations", "ping.target")
	}
	return allSet, missingConfig
}

func readInConfig(v *viper.Viper) {
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

/*
func main() {
	initConfig()
	updateConfig()
	fmt.Println(viper.IsSet("ping.target"))
	for _, element := range pingTarget {
		fmt.Println(element)
	}
}*/
