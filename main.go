package main

import (
	"bufio"
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type checkResult struct {
	success bool
}

func logconfig() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func main() {
	startTime := time.Now()
	logconfig()

	var (
		filename    = flag.String("f", "", "path to urls files")
		timeout     = flag.Int("t", 10, "timeout in seconds")
		concurrency = flag.Int("c", 10, "number of concurrent requests")
	)
	flag.Parse()

	if *filename == "" {
		slog.Error("missing file argument", "usage", "go run main.go -f <path-to-urls.txt>")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*filename)
	if err != nil {
		slog.Error("failed to open file", "file", *filename, "error", err)
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		slog.Error("failed to read file", "file", *filename, "error", err)
		os.Exit(1)
	}

	slog.Info("monitoring started",
		"total_urls", len(lines),
		"timeout", *timeout,
		"concurrency", *concurrency)

	// Create a context that is canceled on an interrupt signal (Ctrl+C).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := &http.Client{Timeout: time.Duration(*timeout) * time.Second}
	results := make([]checkResult, len(lines))
	var wg sync.WaitGroup

	sem := make(chan struct{}, *concurrency)

	for i, raw := range lines {
		i := i
		orig := raw
		url := raw
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(idx int, u, orig string) {
			defer wg.Done()
			defer func() { <-sem }()

			statusCode, status, err := CheckOnce(ctx, client, orig)

			if err != nil {
				results[idx] = checkResult{success: false}
				slog.Error("request creation failed",
					"line_number", idx+1,
					"url", orig,
					"error", err)
				return
			}

			if statusCode >= 200 && statusCode < 300 {
				results[idx] = checkResult{success: true}
				slog.Info("http check completed",
					"line_number", idx+1,
					"url", orig,
					"statuscode", statusCode,
					"status", status)
			} else {
				results[idx] = checkResult{success: false}
				slog.Warn("http check failed (bad status)",
					"line_number", idx+1,
					"url", orig,
					"statuscode", statusCode,
					"status", status)
			}
		}(i, url, orig)
	}

	wg.Wait()

	var success, failed int
	for _, r := range results {
		if r.success {
			success++
		} else {
			failed++
		}
	}

	slog.Info("monitoring completed",
		"total_urls", len(lines),
		"timeout", *timeout,
		"concurrency", *concurrency,
		"success", success,
		"failed", failed,
		"duration", time.Since(startTime).String(),
		"timestamp", time.Now().Format(time.RFC3339))
}
