# oreno-mssh

[![v0.0.3](https://img.shields.io/badge/package-v0.0.3-ff69b4.svg)]()
[![GoDoc](https://godoc.org/github.com/kenzo0107/omssh?status.svg)](https://godoc.org/github.com/kenzo0107/omssh)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/kenzo0107/omssh/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kenzo0107/omssh)](https://goreportcard.com/report/github.com/kenzo0107/omssh)
[![codecov](https://codecov.io/gh/kenzo0107/omssh/branch/master/graph/badge.svg)](https://codecov.io/gh/kenzo0107/omssh)

## AWS mssh Tool wrapper

This project provides an interactive cli tool using ec2-instance-connect api to connect via ssh to ec2 instance.

* select profile in your aws credentials
* select ec2 instance id in list of ec2 intances gotten by credentials specified profile
* ssh to ec2 instance by using ec2 instance connect api

## Install  

```
$ go get -u github.com/kenzo0107/omssh
```

## LICENSE

The MIT License (MIT)

Copyright (c) 2019 Kenzo Tanaka
