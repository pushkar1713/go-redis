# go-redis

`go-redis` is a minimal Redis-compatible, in-memory key-value store written in Go.

It currently supports a small subset of RESP commands and is intended as an educational/experimental project.

Repository: [github.com/pushkar1713/go-redis](https://github.com/pushkar1713/go-redis)

## Status

`Experimental (v0.1)`

## Features

- TCP server on port `8080`
- RESP array + bulk-string command parsing
- Concurrent client handling (goroutine per connection)
- In-memory data storage with RW mutex
- Minimal Redis compatibility layer

## Supported Commands

| Command | Description | Example |
| --- | --- | --- |
| `PING` | Health check | `PING` |
| `SET key value` | Set a key | `SET name PUSHKAR` |
| `GET key` | Get a key | `GET name` |
| `DEL key` | Delete a key | `DEL name` |
| `EXISTS key [key ...]` | Count existing keys | `EXISTS name age city` |
| `COMMAND` | Redis compatibility command | `COMMAND` |

## Quick Start

### Prerequisites

- Go `1.25.5` or later
- `redis-cli` (optional, for manual testing)

### Clone

```bash
git clone https://github.com/pushkar1713/go-redis.git
cd go-redis
```

### Build

```bash
go build -o go-redis .
```

### Run

```bash
./go-redis
```

Server listens on:

```text
localhost:8080
```

### Test With `redis-cli`

```bash
redis-cli -p 8080
```

Example session:

```text
127.0.0.1:8080> PING
PONG
127.0.0.1:8080> SET name PUSHKAR
OK
127.0.0.1:8080> GET name
"PUSHKAR"
127.0.0.1:8080> EXISTS name age
(integer) 1
```

## Implementation Notes

- All data is stored in-memory (no persistence).
- This is not production-ready and omits many Redis features.
- Current parser uppercases processed tokens, so values may be stored in uppercase.

## Benchmarks

Benchmarks were run locally with `redis-benchmark`.

Configuration:

- Requests: `100000`
- Clients: `50`
- Payload: `3 bytes`
- Keepalive: `enabled`

### Commands Used

```bash
redis-benchmark -p 8080 -t set -n 100000 -c 50 -r 100000
redis-benchmark -p 8080 -t get -n 100000 -c 50 -r 100000
redis-benchmark -p 8080 -t set,get -n 100000 -c 50 -r 100000
```

### Results Summary

| Workload | go-redis Throughput | go-redis Avg Latency | Redis Throughput | Redis Avg Latency |
| --- | --- | --- | --- | --- |
| `SET` | ~1600 req/sec | ~0.19 ms | ~5050 req/sec | ~9.7 ms |
| `GET` | ~2781 req/sec | ~17 ms | ~3267 req/sec | ~15 ms |
| `SET` (mixed) | ~639 req/sec | - | ~5325 req/sec | - |
| `GET` (mixed) | ~407 req/sec | - | ~8719 req/sec | - |

## Roadmap

- Add graceful error handling for malformed RESP input
- Improve command validation and argument handling
- Add persistence options (AOF/RDB style)
- Add tests and CI
