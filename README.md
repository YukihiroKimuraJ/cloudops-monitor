# CloudOps Monitor

🚧 **Work in Progress** - A lightweight URL monitoring tool written in Go

## Features
- Concurrent URL health checks from a file
- Support for comments and blank lines in URL list
- Auto-adds https:// to URLs without protocol
- Graceful shutdown with Ctrl+C
- Configurable timeout (10s) and concurrency (10 workers)

## Usage
```bash
go run main.go urls.txt
```

## Example urls.txt
```
https://example.com
https://google.com
# This is a comment
example.org  # Auto-adds https://
```

## Roadmap
- [ ] Add Prometheus metrics export
- [ ] Implement retry logic with exponential backoff
- [ ] Support custom HTTP headers
- [ ] Add structured logging
- [ ] Create configuration file support

## Requirements
- Go 1.21 or later

## License
MIT
