# CloudOps Monitor

[![CI](https://github.com/YukihiroKimuraJ/cloudops-monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/YukihiroKimuraJ/cloudops-monitor/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/YukihiroKimuraJ/cloudops-monitor)](https://goreportcard.com/report/github.com/YukihiroKimuraJ/cloudops-monitor)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, concurrent HTTP endpoint monitoring tool written in Go.

## Features

- 🚀 **Concurrent Health Checks** - Monitor multiple URLs simultaneously with configurable worker pool size
- 📊 **Structured Logging** - JSON-formatted logs using Go's `slog` package for easy integration with log aggregators
- ⚙️ **Flexible Configuration** - Command-line flags for timeout, concurrency, and input file
- 🛑 **Graceful Shutdown** - Clean termination with Ctrl+C signal handling
- 🔗 **Smart URL Handling** - Automatically prepends https:// to URLs without a scheme, supports comments in URL files
- ✅ **Well Tested** - Comprehensive unit tests including timeout handling
- 🐳 **Docker Support** - Containerized deployment with multi-stage build

## Quick Start

### Installation

```bash
git clone https://github.com/YukihiroKimuraJ/cloudops-monitor.git
cd cloudops-monitor
go build -o cloudops-monitor .
```

### Basic Usage

```bash
# Create a URL list file
cat <<EOF > urls.txt
https://google.com
https://github.com
example.com
EOF

# Run the monitor
./cloudops-monitor -f urls.txt
```

## Command-Line Options

| Flag | Description                              | Default |
|------|------------------------------------------|---------|
| `-f` | Path to URL list file (required)         | -       |
| `-t` | HTTP request timeout (seconds)           | `10`    |
| `-c` | Number of concurrent workers             | `10`    |

### Examples

```bash
# Monitor with 30 second timeout
./cloudops-monitor -f urls.txt -t 30

# High concurrency (50 workers)
./cloudops-monitor -f urls.txt -c 50

# Example: production-like settings
./cloudops-monitor -f production-urls.txt -t 15 -c 20
```

## Docker

```bash
# Build the Image
docker build -t cloudops-monitor .

# Run with volume mount (recommended)
docker run --rm -v $(pwd)/urls.txt:/app/urls.txt cloudops-monitor -f /app/urls.txt

# With custom timeout and concurrency
docker run --rm -v $(pwd)/urls.txt:/app/urls.txt cloudops-monitor -f /app/urls.txt -t 30 -c 20
```

### Docker Image Details

- Base Image: Alpine Linux (lightweight)
- Build Method: Multi-stage build for optimized image size
- Security: Runs as non-root user
- URL Files: Must be mounted as volume (not included in image)

## URL File Format

```text
# Production services
https://api.example.com/health
https://web.example.com

# External dependencies
google.com          # Protocol auto-added
github.com/status

# This line is a comment
```

- Lines starting with `#` are treated as comments
- Empty lines are ignored
- URLs without protocol automatically get `https://` prefix

## Output

CloudOps Monitor outputs structured JSON logs for easy parsing:

```json
{"time":"2025-12-21T09:32:06.530644+09:00","level":"INFO","msg":"monitoring started","total_urls":5,"timeout":20,"concurrency":20}
{"time":"2025-12-21T09:32:06.574340+09:00","level":"INFO","msg":"http check completed","line_number":2,"url":"https://github.com","statuscode":200,"status":"200 OK"}
{"time":"2025-12-21T09:32:06.892173+09:00","level":"INFO","msg":"monitoring completed","total_urls":5,"success":5,"failed":0,"duration":"361.529ms"}
```

### Log Levels

| Level | Description |
|-------|-------------|
| `INFO` | Successful checks and monitoring status |
| `WARN` | Non-2xx status codes (server responded but with error) |
| `ERROR` | Connection failures, timeouts, invalid URLs |

## Development

### Prerequisites

- Go 1.22 or later

### Project Structure

```
cloudops-monitor/
├── main.go           # Entry point, CLI, and orchestration
├── monitor.go        # Core monitoring logic (CheckOnce, normalizeURL)
├── monitor_test.go   # Unit tests
├── Dockerfile        # Docker image definition
├── .dockerignore     # Docker build exclusions
├── urls.txt          # Sample URL list
├── .github/
│   └── workflows/
│       └── ci.yml    # GitHub Actions CI
└── README.md
```

### Running Tests

```bash
# Run all tests
go test -v

# Run with coverage
go test -v -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Cases

| Test | Description |
|------|-------------|
| `TestCheckOnce_Status200` | Verifies successful HTTP 200 response handling |
| `TestCheckOnce_Status500` | Verifies HTTP 500 error response handling |
| `TestCheckOnce_Timeout` | Verifies timeout error detection |
| `TestNormalizeURL` | Verifies URL normalization (table-driven test) |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  main.go                                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ CLI Parsing │→ │ File Reader │→ │ Concurrent Workers  │  │
│  └─────────────┘  └─────────────┘  └──────────┬──────────┘  │
│                                               │             │
│                                    ┌──────────▼──────────┐  │
│                                    │ monitor.go          │  │
│                                    │ - CheckOnce()       │  │
│                                    │ - normalizeURL()    │  │
│                                    └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Key Components

| Component | File | Responsibility |
|-----------|------|----------------|
| CLI & Orchestration | `main.go` | Flag parsing, file reading, worker pool management |
| Signal Handling | `main.go` | Graceful shutdown on Ctrl+C |
| HTTP Client | `monitor.go` | Single URL check, URL normalization |

## Roadmap

### High Priority

- [x] Docker support
- [ ] Retry logic with exponential backoff

### Medium Priority

- [ ] Prometheus metrics endpoint (`/metrics`)
- [ ] YAML/TOML configuration file

### Future Ideas

- [ ] Alerting integrations (Slack, PagerDuty, Email)
- [ ] Response time tracking and statistics

## Blog Post

開発の背景や学んだことについて詳しくはこちら：

📝 [レガシーインフラエンジニアがGo言語でURL監視ツールを作った話](https://zenn.dev/yukihirokimuraj/articles/13ee8236b029d6) (Japanese)

📝 [Go 1.21 slogで構造化ログを実装する](https://zenn.dev/yukihirokimuraj/articles/2a79afcbaaec05) (Japanese)

## Author

Yukihiro Kimura

- Infrastructure Engineer with 7+ years of experience
- AWS All 12 Certifications holder
- Currently learning Go and transitioning to modern cloud-native infrastructure

[![GitHub](https://img.shields.io/badge/GitHub-YukihiroKimuraJ-181717?logo=github)](https://github.com/YukihiroKimuraJ)
[![Zenn](https://img.shields.io/badge/Zenn-yukihirokimuraj-3EA8FF?logo=zenn)](https://zenn.dev/yukihirokimuraj)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
