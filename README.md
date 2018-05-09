# CGW-Exporter

This exporter provides metrics for the *CGW* service in the *prometheus* compatible format.

It is orginally based on [ping_exporter](https://github.com/czerwonk/ping_exporter) by Daniel Czerwonk.

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
  
