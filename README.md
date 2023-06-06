## po-test - prometheus operator test CLI
`po-test` allows you to run [prometheus alert unit tests](https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/) against
prometheus operator manifests. Simply point your unit tests at a prometheus operator `PrometheusRule` E.g.

```
rule_files:
  - prometheus-operator-rules.yml
```

Where `prometheus-operator-rules.yml` looks something like:
```
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    prometheus: k8s
    role: alert-rules
  name: prometheus-rules
spec:
  groups:
    - name: alerts
...
```

Then run `po-test prometheus-operator-unittest.yml`. 
The CLI will strip the kubernetes metadata from the target, run the unit tests against it, then put it back again. 

It is also suitable for running as part of a CI/CD pipeline.

## Prerequisites

### Promtool
Run this to install prometheus that will also provide promtool
```brew install prometheus```

## Install

```
# go >= 1.17
# Using `go get` to install binaries is deprecated.
# The version suffix is mandatory.
go install github.com/loveholidays/po-test@latest

# go < 1.17
go get github.com/loveholidays/po-test
```

### Homebrew
```
brew install loveholidays/taps/po-test
```

### Linux
Grab the latest OS/Arch compatible binary from our [Releases](https://github.com/loveholidays/po-test/releases) page.

### From source
```bash
git clone git@github.com:loveholidays/po-test.git
make build
```

## Usage
```
po-test test-file1.yaml test-file2.yaml ...
```
