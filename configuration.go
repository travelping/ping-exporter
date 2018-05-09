package main

import (
	"github.com/spf13/viper"
	//flag "github.com/spf13/pflag"
	"fmt"
	"strings"
	"time"
)

var (
	//metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which metrics will be exposed")
	//listenAddress = flag.String("web.listen-address", ":9427", "Address on which to expose metrics and web interface")
	//pingInterval  = flag.Duration("ping.interval", time.Duration(5)*time.Second, "Interval for ICMP echo requests")
	//pingTimeout   = flag.Duration("ping.timeout", time.Duration(4)*time.Second, "Timeout for ICMP echo request")
	//pingTarget    = flag.String("ping.target", "9.9.9.9", "IP address of target")
	//pingSourceV4  = flag.String("ping.source.ipv4", "0.0.0.0", "Source Address of ICMP echo requests")
	//pingSourceV6  = flag.String("ping.source.ipv6", "::", "Source Address of ICMP echo requests")
	//dnsRefresh    = flag.Duration("dns.refresh", time.Duration(1)*time.Minute, "Interval for refreshing DNS records and updating targets accordingly (0 if disabled)")

	metricsPath   string
	listenAddress string
	pingInterval  time.Duration
	pingTimeout   time.Duration
	pingTarget    []string
	pingSourceV4  string
	pingSourceV6  string
	dnsRefresh    time.Duration
)

var (
	mandatoryConfig = [...]string{"web.listen-address", "ping.interval", "ping.timeout", "ping.source.ipv4", "ping.source.ipv6", "dns.refresh", "ping.target", "web.telemetry-path"}
	optionalConfig  = [...]string{}
)

func initConfig() {
	setDefaults()
	viper.SetConfigName("cgw-exporter")
	viper.AddConfigPath("/etc/defaults/cgw-exporter/")
	viper.AddConfigPath("/etc/cgw-exporter/")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("CGWEXPORTER")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	bindEnvVariables()
}

func setDefaults() {
	viper.SetDefault("web.listen-address", ":9427")
	viper.SetDefault("web.telemetry-path", "/metrics")
	viper.SetDefault("ping.interval", "5s")
	viper.SetDefault("ping.timeout", "4s")
	viper.SetDefault("ping.source.ipv4", "0.0.0.0")
	viper.SetDefault("ping.source.ipv6", "::")
	viper.SetDefault("dns.refresh", "1m")
}

func updateConfig() {
	listenAddress = viper.GetString("web.listen-address")
	metricsPath = viper.GetString("web.telemetry-path")
	pingInterval = viper.GetDuration("ping.interval")
	pingTimeout = viper.GetDuration("ping.interval")
	pingTarget = viper.GetStringSlice("ping.target")
	pingSourceV4 = viper.GetString("ping.source.ipv4")
	pingSourceV6 = viper.GetString("ping.source.ipv6")
	dnsRefresh = viper.GetDuration("dns.refresh")
}

func bindEnvVariables() {
	//flag.Parse()
	//viper.BindPFlags(flag.CommandLine)

	for _, element := range mandatoryConfig {
		viper.BindEnv(element)
	}
	for _, element := range optionalConfig {
		viper.BindEnv(element)
	}
}

func isMandatoryConfigSet() (bool, []string) {
	allSet := true
	var missingConfig []string
	for _, element := range mandatoryConfig {
		if !viper.IsSet(element) {
			allSet = false
			missingConfig = append(missingConfig, element)
		}
	}
	return allSet, missingConfig
}

func readInConfig() {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
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
