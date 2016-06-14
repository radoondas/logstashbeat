[![Build Status](https://travis-ci.org/radoondas/logstashbeat.svg?branch=master)](https://travis-ci.org/radoondas/logstashbeat)
[![GoReportCard](https://goreportcard.com/badge/github.com/radoondas/logstashbeat)](https://goreportcard.com/report/github.com/radoondas/logstashbeat)

# Logstashbeat

**There is another [logstashbeat](https://github.com/consulthys/logstashbeat) which was developed in parallel without any knowledge of each other. I worked on this one for personal purposes. I do not want to compete but I don't want to delete it as well. :)**

**Note: This is development version**

**Note: Beat works only with Logstash 5.0aplha3**

Welcome to Logstashbeat which queries metrics from Logstash and index them in Elasticsearch.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/radoondas`

## Getting Started with Logstashbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Build

To build the binary for Logstashbeat run the command below. This will generate a binary
in the same directory with the name logstashbeat.

```
make
```


### Run

To run Logstashbeat with debugging output enabled, run:

```
./logstashbeat -c logstashbeat.yml -e -d "*"
```


### Test (not implemented yet)

To test Logstashbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


### Package

To be able to package Logstashbeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Logstashbeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/logstashbeat.template.json and etc/logstashbeat.asciidoc

```
make update
```


### Cleanup

To clean  Logstashbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Logstashbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/radoondas
cd ${GOPATH}/github.com/radoondas
git clone https://github.com/radoondas/logstashbeat
```
