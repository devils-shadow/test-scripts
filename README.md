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

This repository includes a minimal stub of `github.com/livekit/server-sdk-go`
under `stubs/`. The `go.mod` file uses a `replace` directive to point to this
stub so the project can compile without internet access.

To perform live testing with the real SDK, remove the `replace` line from
`go.mod` and run `go mod download` to fetch the actual dependency.

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
