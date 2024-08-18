# ibt

    IRacing Telemetry parsing and processing library


<!-- [![](https://img.shields.io/github/actions/workflow/status/spf13/cobra/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/spf13/cobra/actions?query=workflow%3ATest) -->
[![Go Reference](https://pkg.go.dev/badge/github.com/teamjorge/ibt.svg)](https://pkg.go.dev/github.com/teamjorge/ibt)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamjorge/ibt)](https://goreportcard.com/report/github.com/teamjorge/ibt)



## Install

`go get github.com/teamjorge/ibt`


## Overview

`ibt` is a package from parsing IRacing telemetry files. An *ibt* file is created when you enter the car and ends when you exit the car. By default, you can find these files in your `IRacing/telemetry/[car]/` directory. These files are binary for the most part, with the exception of the session data.

This package does not (yet) parse real-time as that requires opening a memory-mapped file and CGO. A planned real-time parsing package utilising this one will be available at some point.

## Features

* Easy to use telemetry tick processing interface.
* Quick parsing of file metadata.
* Grouping of *ibt* files into the sessions where they originate from.
* Great test coverage and code documentation.
* Freedom to use it your own way. Most of what is needed are public functions/methods.

