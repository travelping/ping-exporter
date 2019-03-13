# Ping-Exporter

This exporter provides metrics for ICMP Echo requests in the *prometheus*
compatible format.

## Copyright

It is orginally based on
[ping_exporter](https://github.com/czerwonk/ping_exporter) by Daniel Czerwonk.

## Configuration
### Configuration formats

The configuration can either be set using `yaml` as a format or environmental
variables.

For configurations in yaml like:

```yaml
layer1:
  layer2:
    layer3: true
```

The value key would be called `layer1.layer2.layer3` with the value `true`
when mentioned in this documentation. To use *environmental variables* for
configuration, the key translates to `PINGEXPORTER_LAYER1_LAYER2_LAYER3`.
So every `.` will be replaced by `_` and the name will be all caps.

Further list items will be seperated by space, e.g.
`PINGEXPORTER_TARGET="192.0.2.1 192.0.2.2"`

The configuration in `yaml` format has to be saved under
`/etc/ping-exporter/ping-exporter.yaml` as of now.
The priority of environmental variables is higher than the one of the
configuration file and therefore will override the values.

### Configuration parameters

The following configuration parameters are available so far:

```yaml
version: "1.0" # version of the configuration (mandatory)
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
  
### Multiple ping configurations

As of version 0.2.0 it is also possible to use multiple configurations for the
ping targets. One use case is use multiple source IP addresses.

When `ping.configurations` is set, `ping.target`, `ping.source`,
`ping.interval` and `ping.timeout` will not be evaluated and consequently
also the corrensponding command line values will be ignored.

Example configuration:

```yaml
version: "1.0"  # version of the configuration (mandatory)
ping:
  configurations:
    - sourceV4: 192.0.2.1          # Source address of ICMP requests
      sourceV6: "2001:0DB8:1::1"   # Source address of ICMP requests
      pingInterval: 5s             # interval for ICMP requests
      pingTimeout: 4s              # timeout for ICMP requests
      pingTargets:                 # list of ICMP targets
        - pingTarget: 192.0.2.10
        - pingTarget: 198.51.100.1
    - sourceV4: 192.0.2.2
      sourceV6: "2001:0DB8:2::2"
      pingInterval: 5s
      pingTimeout: 4s
      pingTargets:
        - pingTarget: 203.0.113.1
        - pingTarget: "2001:0DB8:2::10"
```

#### Using additional labels

Is is possible to add additional labels to the configuration.

The labels (both source and target) will be added to the resulting metrics.
Therefore the label names have to be different between source and target.

This is a working configuration:

```yaml
version: "1.0"  # version of the configuration (mandatory)
ping:
  configurations:
    - sourceV4: 192.0.2.1          # Source address of ICMP requests
      sourceV6: "2001:0DB8:1::1"   # Source address of ICMP requests
      sourceLabels:
        source_name: server1       # a source label
        source_company: tier1prov  # another source label
      pingInterval: 5s             # interval for ICMP requests
      pingTimeout: 4s              # timeout for ICMP requests
      pingTargets:                 # list of ICMP targets
        - pingTarget: 192.0.2.10
          targetLabels:
            target_name: server11    # a target label
        - pingTarget: 198.51.100.1
          targetLabels:
            target_name: server12    # another target label
```


Attention!: the target and source label names have to be the same for all configuration, as it will otherwise result in a runtime error.

Even in two different confiugurations with different source IP address, the labels for source and targets have to match.

For example the following is NOT valid:

```yaml
version: "1.0"  # version of the configuration (mandatory)
ping:
  configurations:
    - sourceV4: 192.0.2.1          # Source address of ICMP requests
      sourceV6: "2001:0DB8:1::1"   # Source address of ICMP requests
      sourceLabels:
        source_name: server1       # a source label
        source_company: tier1prov  # another source label
      pingInterval: 5s             # interval for ICMP requests
      pingTimeout: 4s              # timeout for ICMP requests
      pingTargets:                 # list of ICMP targets
        - pingTarget: 192.0.2.10
          targetLabels:
            target_name: server11    # a target label
        - pingTarget: 198.51.100.1
          targetLabels:
            target_company: server12    # NOT VALID configuration 
```
