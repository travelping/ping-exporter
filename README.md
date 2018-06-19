# CGW-Exporter

This exporter provides metrics for the *CGW* service in the *prometheus* compatible format.

## copyright

It is orginally based on [ping_exporter](https://github.com/czerwonk/ping_exporter) by Daniel Czerwonk.

> Copyright (c) 2018 Daniel Czerwonk
> 
> Copyright (c) 2018 Travelping GmbH


## configuration
### configuration formats

The configuration can either be set using `yaml` as a format or environmental variables.

For configurations in yaml like:

```yaml
layer1:
  layer2:
    layer3: true
```

The value key would be called `layer1.layer2.layer3` with the value `true` when mentioned in this documentation.
To use *environmental variables* for configuration, the key translates to `CGWEXPORTER_LAYER1_LAYER2_LAYER3`.
So every `.` will be replaced by `_` and the name will be all caps.

Further list items will be seperated by space, e.g. `CGW_EXPORTER_TARGET="192.0.2.1 192.0.2.2"`

The configuration in `yaml` format has to be saved under `/etc/cgw-exporter/cgw-exporter.yaml` as of now.
The priority of environmental variables is higher than the one of the configuration file and 
therefore will override the values.

### configuration parameters

The following configuration parameters are available so far:

```yaml
web:
  listen-address: ":9427" # binding, the http server will be listening on
  telemetry-path: "/metrics" # path, under which the metrics will be exposed
ping:
  interval: "5s" # interval for ICMP requests
  timeout: "4s"  # timeout for ICMP requests
  source:
    ipv4: "0.0.0.0" # Source address of ICMP requests
    ipv6: "::"      # Source address of ICMP requests
  target:
  - 192.0.2.1
  - 192.0.2.2
dns:
  refresh: "1m" # Interval for refreshing DNS records and updating targets accordingly (0 if disabled)
```
  
### multiple ping configurations

As of version 0.2.0 it is also possible to use multiple configurations for the ping targets.
One use case is use multiple source IP addresses.

When `ping.configurations` is set, `ping.target`, `ping.source`, `ping.interval` and `ping.timeout` will not be evaluated and consequently also the corrensponding command line values will be ignored.

Example configuration:

```yaml
ping:
  configurations:
    - sourceV4: 192.0.2.1          # Source address of ICMP requests
      sourceV6: "2001:0DB8:1::1"   # Source address of ICMP requests
      pingInterval: 5s             # interval for ICMP requests
      pingTimeout: 4s              # timeout for ICMP requests
      pingTargets:                 # list of ICMP targets
        - 192.0.2.10
        - 198.51.100.1
    - sourceV4: 192.0.2.2
      sourceV6: "2001:0DB8:2::2"
      pingInterval: 5s
      pingTimeout: 4s
      pingTargets:
        - 203.0.113.1
        - "2001:0DB8:2::10"
```
