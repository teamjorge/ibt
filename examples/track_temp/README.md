# track_temp

The `track_temp` example shows a simple processor summarising the track temperature on each lap of the provided `ibt` files.

## Running

From the root of the repository:

```shell
go run examples/track_temp/main.go

# Or with your own files

go run examples/track_temp/main.go /path/to/telem/files/*.ibt
```

Using the included `ibt` file will yield only a single lap and it's value. However, if you have telemetry consisting of a few laps and/or files from a longer session, you should have a nicely summarised per-lap output.