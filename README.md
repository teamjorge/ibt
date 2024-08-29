# ibt

    iRacing Telemetry parsing and processing library

![ibt logo](assets/ibt_mascot.png)

[![Go Reference](https://pkg.go.dev/badge/github.com/teamjorge/ibt.svg)](https://pkg.go.dev/github.com/teamjorge/ibt)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamjorge/ibt)](https://goreportcard.com/report/github.com/teamjorge/ibt)
[![codecov](https://codecov.io/gh/teamjorge/ibt/branch/main/graph/badge.svg?token=08QVKSEPXT)](https://codecov.io/gh/teamjorge/ibt)

## Install

`go get github.com/teamjorge/ibt`

## Overview

`ibt` is a package from parsing iRacing telemetry files. An *ibt* file is created when you enter the car and ends when you exit the car. By default, you can find these files in your `iRacing/telemetry/[car]/` directory. These files are binary for the most part, with the exception of the session data.

This package will not parse real-time telemetry as that requires opening a memory-mapped file and CGO. A planned real-time parsing package leverage `ibt` will be available in the future.

## Features

* Easy to use telemetry tick processing interface.
* Quick parsing of file metadata.
* Grouping of *ibt* files into the sessions where they originate from.
* Great test coverage and code documentation.
* Freedom to use it your own way. Most functions/methods has been made public.

## Examples

The [Examples](https://github.com/teamjorge/ibt/tree/main/examples) directory houses all of the available examples.

To try these examples locally, please clone to repository:

```shell
git clone https://github.com/teamjorge/ibt
#or
git clone git@github.com:teamjorge/ibt.git
```

Please have a look at the instructions in the examples [`README`](./examples/README.md) for details on how to run each example.
