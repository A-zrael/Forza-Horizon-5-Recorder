![](https://github.com/A-zrael/Forza-Horizon-5-Recorder/blob/main/Telemetry%20Recorder%20Icon.jpeg?raw=true)

# Forza Horizon 5 Telemetry Recorder (Go)

Lightweight Go UDP listener that records Forza Horizon 5 “Data Out” telemetry to per-car CSV files. Supports multiple simultaneous senders by listening on a configurable list or range of ports.

## Prerequisites

- Go toolchain (go.mod targets 1.25.4).
- Forza Horizon 5 with “Data Out” enabled and pointed to the machine running this app. Use a unique port per console/PC if multiple cars are recording at once.

## Quick start

```bash
# Run directly (default ports 5030-5040)
go run .

# Custom ports (comma list or inclusive range)
go run . -ports 3050,3051
go run . -ports 3050-3040
```

### Build a binary

```bash
go build -o forza_recorder
./forza_recorder -ports 5030-5040
```

## Runtime behavior

- Opens a UDP listener on each configured port and tags incoming packets with the port as the car ID.
- Waits until every seen car reports `IsRaceOn == true`, then begins recording.
- Writes one CSV per car in the working directory named `Car-<port>.csv`.
- Stops once all cars send a packet with `IsRaceOn == false`, then exports the files.

## CLI flags

- `-ports` (string): Comma-separated list or inclusive range (`start-end`) of UDP ports to listen on. Default: `5030-5040`.

## Output schema

CSV headers match the fields in `models.Carstate` (timestamp, RPMs, acceleration/velocity, suspension, tire slip/temps, position, speed, lap info, controls, etc.). See `models/models.go` for the full column list.

## Race viewer

The original in-browser race viewer is no longer part of this repo. If you want visualization/analysis, use the separate viewer project (not included here) or import the generated CSVs into your own tooling (Excel, pandas, etc.).
