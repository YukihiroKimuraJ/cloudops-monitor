package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func run(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("run: failed to open file %q: %w", file, err)
	}
	// ... work, return errors rather than calling os.Exit
	return f, nil
}

func main() {
	var f *os.File
	var err error

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error : A file path argument is required.")
		fmt.Fprintln(os.Stderr, "Usage : go run main.go <path-to urls.txt")
		os.Exit(1)
	}
	file := os.Args[1]
	f, err = run(file)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "read %s: %v\n", file, err)
		os.Exit(1)
	}

	// Create a context that is canceled on an interrupt signal (Ctrl+C).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := &http.Client{Timeout: 10 * time.Second}
	results := make([]string, len(lines))
	var wg sync.WaitGroup

	const concurrency = 10
	sem := make(chan struct{}, concurrency)

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

			req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
			if err != nil {
				results[idx] = fmt.Sprintf("%s -> ERROR: %v", orig, err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				results[idx] = fmt.Sprintf("%s -> ERROR: %v", orig, err)
				return
			}
			defer resp.Body.Close()
			results[idx] = fmt.Sprintf("%s -> %d %s", orig, resp.StatusCode, resp.Status)
		}(i, url, orig)

	}

	wg.Wait()

	for _, r := range results {
		fmt.Println(r)
	}

}

