# LiveKit Load Test Scripts

This repository contains a small load testing tool written in Go and a shell
helper for launching multiple bots. It can be used to spawn rooms of fake
participants against a LiveKit server for basic stress testing.

## Setup

1. Install **Go 1.24.3** on your machine.
2. Clone this repository and change into the project directory.
3. Download `sample720p.ivf` and `sample.ogg` from the [LiveKit example
   media folder](https://github.com/livekit/client-sdk-js/tree/main/examples/media)
   and place them in this directory. These files are published by each bot when
   it joins a room.
   If you have `curl` installed you can also run `./get-media.sh` to fetch them
   automatically.

## Stubbed dependency

The project now depends on the real
`github.com/livekit/server-sdk-go/v2` module in `go.mod`. A minimal stub of
the SDK is still provided under `stubs/` for situations where you need to build
without network access.

To use the stub, add the following line to `go.mod`:

```
replace github.com/livekit/server-sdk-go/v2 => ./stubs/github.com/livekit/server-sdk-go/v2
```

After adding the replace directive, run `go mod tidy` so Go picks up the local
version. Remove the line again when you want to fetch the real dependency.

## Setup script

Before running the load tests you can execute `./setup.sh` to fetch the Go
dependencies and download the sample media files:

```bash
./setup.sh
```

## Running the load test

The easiest way to run the test is via the interactive helper script:

```bash
./run-custom.sh
```

The script will prompt for the number of rooms, bots per room and duration of
the test. It then calls `go run loadbot.go` with the selected parameters.

Alternatively you can run the Go program directly:

```bash
go run loadbot.go -rooms <num> -bots <num> -d <duration>
```

See `go run loadbot.go -h` for all available flags.

## Optional logging

Passing the `-log` flag (or answering `y` when prompted by the script) will
write NDJSON formatted logs to `last_run.json`. Each line describes a bot event
(join, leave, errors, etc.) with a timestamp, making it easier to inspect or
process later.
