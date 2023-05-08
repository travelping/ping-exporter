# Changelog

## pre v1.0.0

### 0.7.0

* update libraries to newest versions
* Use Alpine 3.17 as base image

### 0.6.0

* [breaking] add labels for source and target
  * labels for the source and destinations can be set
    to be added in the resulting metrics
  * the configuration for multiple sources has therefore changed
* [breaking] start versioning of configuration file
  * the configuration has to contain a version number to ensure backward compatibility
    in the future

### 0.5.1

* fix problem with missing configuration file
  * if the default configuration file is used, it is assumed, that
    configuration with environmental variables and command line
    parameters is still possible.
  * now if the default configuration file is missing, it will not throw
    an error but just use env vars and cli parameters.
  * if a config file path is given via cli parameter and is missing
    the program will exit with an error.

### 0.5.0

* CLI: provide -h, -v and -c <config-file>
* CLI: remove unused flags --uniform.domain, --norm.domain, --normal.mean
* Docker: contain org.label-schema labels
* Docker: use Alpine-3.8
* Docker: do not create an empty `/etc/ping-exporter/ping-exporter.yaml`

### 0.4.0

* remove `/etc/defaults/ping-exporter/` search path, because just one
  is used which should be `/etc/ping-exporter/`.

### 0.3.1

* clean up help text

### 0.3.0

* move to Github and change name from *cgw-exporter* to *ping-exporter* including
  configuration prefixes and paths

### 0.2.1

* fix error in configuration parsing

### 0.2.0

Add feature of configuring multiple source and target address combinations for
ICMP metrics using yaml configuration files.

### 0.1.0

Initial working prototype of *cgw-exporter* with feature of ICMP metrics with
single source address and multiple target addresses
