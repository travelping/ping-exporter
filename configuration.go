package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/prometheus/common/log"
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
	pingConfigurations []pingConfig
	dnsRefresh         time.Duration
}

var (
	mandatoryConfig = [...]string{"web.listen-address", "ping.interval", "ping.timeout", "ping.source.ipv4", "ping.source.ipv6", "dns.refresh", "web.telemetry-path"}
	optionalConfig  = [...]string{"ping.target"}
)

type pingConfig struct {
	SourceV4     string             // source address of ICMP requests
	SourceV6     string             // source address of ICMP requests
	SourceLabels map[string]string  //labels for source addresses
	PingTargets  []pingTargetConfig // target addresses of ICMP requests
	PingInterval time.Duration      // interval between ICMP requests
	PingTimeout  time.Duration      // timeout of ICMP requests
}

type pingTargetConfig struct {
	PingTarget   string            // destination IP address
	TargetLabels map[string]string // labels for for targets to export to metrics
}

func initConfig(v *viper.Viper, flags *pflag.FlagSet) {
	setDefaults(v)

	v.BindPFlags(flags)

	v.SetEnvPrefix("PINGEXPORTER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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

func newConfiguration(flags *pflag.FlagSet) (*Configuration, error) {
	v := viper.GetViper()
	initConfig(v, flags)

	cfile, _ := flags.GetString("config")
	if cfile != "" {
		if _, err := os.Stat(cfile); !os.IsNotExist(err) {
			v.SetConfigFile(cfile)
			if err := v.ReadInConfig(); err != nil {
				return nil, err
			}
		}
	}

	config, err := (new(Configuration)).updateConfig(v)
	if err != nil {
		return nil, err
	}
	if valid, missing := isMandatoryConfigSet(v); !valid {
		err = fmt.Errorf("configuration parameters", missing, "missing")
	}
	return config, err
}

func (conf *Configuration) updateConfig(v *viper.Viper) (*Configuration, error) {
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
			return nil, err
		}
	} else {
		conf.hasPingMultiConfig = false
	}
	return conf, nil
}

func bindEnvVariables(v *viper.Viper) {
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
