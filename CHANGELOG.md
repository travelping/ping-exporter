# Changelog

## pre v1.0.0

### 0.4.0

* remove `/etc/defaults/ping-exporter/` search path, because just one
  is used which should be `/etc/ping-exporter/`.

### 0.3.1

* clean up help text

### 0.3.0

* move to github and change name from cgw-exporter to ping-exporter including configuration prefixes and paths

### 0.2.1

* fix error in configuration parsing

### 0.2.0
Add feature of configuring multiple source and target address combinations for ICMP metrics
using yaml configuration files.

### 0.1.0
Initial working prototype of *cgw-exporter* with feature of icmp metrics with single source address
and multiple target addresses
