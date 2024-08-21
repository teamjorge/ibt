# loader

## Overview

The `loader` example shows a pattern for loading telemetry data to an external destination. This could be loading it into a database, posting to an API, or even storing it in a file.

The practices in this example should not be copied verbatim, but should show how `ibt` can be used to achieve this.

Goals of the example:

* Parse the `ibt` files into groups
* For each group of `ibts` process each tick of telemetry data
* When the threshold of processed telemetry ticks have been reached, perform a bulk load to the storage client
* Add a group number to each telemetry tick to ensure they are easily filtered in the external storage layer
* Store the number of batches loaded and print it after processing each group

## Running

```shell
go run examples/track_temp/*.go

# Or with your own files

go run examples/track_temp/*.go /path/to/telem/files/*.ibt
```