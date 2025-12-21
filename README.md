# CloudOps Monitor

A lightweight, production-ready URL monitoring tool written in Go.

## Features

- **Concurrent Health Checks** - Monitor multiple URLs simultaneously with configurable worker pools
- **Structured Logging** - JSON-formatted logs using Go's `slog` package for easy integration with log aggregators
- **Flexible Configuration** - Command-line flags for timeout, concurrency, and input file
- **Graceful Shutdown** - Clean termination with Ctrl+C signal handling
- **Smart URL Handling** - Auto-adds `https://` to URLs without protocol, supports comments in URL files

## Installation

```bash
git clone https://github.com/YukihiroKimuraJ/cloudops-monitor.git
cd cloudops-monitor
go build -o cloudops-monitor main.go
```

## Usage

### Basic Usage

```bash
./cloudops-monitor -f urls.txt
```

### Command-Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-f` | Path to URL list file | `urls.txt` |
| `-t` | HTTP request timeout (seconds) | `10` |
| `-c` | Number of concurrent workers | `10` |

### Examples

```bash
# Monitor URLs with custom timeout
./cloudops-monitor -f urls.txt -t 30

# High concurrency monitoring
./cloudops-monitor -f urls.txt -c 50

# Full options
./cloudops-monitor -f production-urls.txt -t 15 -c 20
```

## URL File Format

```text
# Production services
https://api.example.com/health
https://web.example.com

# External dependencies
google.com          # Protocol auto-added
github.com/status
```

- Lines starting with `#` are treated as comments
- Empty lines are ignored
- URLs without protocol automatically get `https://` prefix

## Output

CloudOps Monitor uses structured logging (JSON format) for machine-readable output:

```json
{"time":"2025-12-21T09:32:06.530644+09:00","level":"INFO","msg":"monitoring started","total_urls":5,"timeout":20,"concurrency":20}
{"time":"2025-12-21T09:32:06.57434+09:00","level":"INFO","msg":"http check completed","line_number":2,"url":"https://github.com","statuscode":200,"status":"200 OK"}
```

## Roadmap

- [ ] Prometheus metrics endpoint (`/metrics`)
- [ ] Retry logic with exponential backoff
- [ ] Custom HTTP headers support
- [ ] YAML/TOML configuration file
- [ ] Alerting integrations (Slack, PagerDuty)

## Requirements

- Go 1.21 or later

## Blog Post

開発の背景や学んだことについて詳しくはこちら：

📝 [レガシーインフラエンジニアがGo言語でURL監視ツールを作った話](https://zenn.dev/yukihirokimuraj/articles/13ee8236b029d6) (Japanese)

## Author

**Yukihiro Kimura**
- Infrastructure Engineer with 8+ years of experience
- AWS All 12 Certifications holder
- Transitioning to modern cloud-native infrastructure

## License

MIT
