# LiveKit Load Test Scripts

This repository contains a small load testing tool written in Go and a shell
helper for launching multiple bots. It can be used to spawn rooms of fake
participants against a LiveKit server for basic stress testing.

## Setup

1. Install **Go 1.24.3** on your machine.
2. Clone this repository and change into the project directory.
3. Ensure [Git LFS](https://git-lfs.com/) is installed so the sample media
   files (`sample720p.ivf` and `sample.ogg`) are downloaded automatically when
   cloning. If they are missing, you can grab them from the [LiveKit example
   media folder](https://github.com/livekit/client-sdk-js/tree/main/examples/media)
  and place them in this directory. These files are published by each bot when
  it joins a room.

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
