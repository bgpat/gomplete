// +build integration

package test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

// Timeout is the timeout seconds.
const Timeout = 5 * time.Second

// WriteAndWait writes the string and waits the terminal response.
func WriteAndWait(ctx context.Context, tb testing.TB, r io.ReadWriter, s string) {
	tb.Helper()
	if _, err := io.WriteString(r, s); err != nil {
		tb.Fatal(err)
	}
	WaitString(ctx, tb, r, s)
}

// WaitString waits the terminal still the presented string is printed.
func WaitString(ctx context.Context, tb testing.TB, r io.Reader, s string) (text string) {
	tb.Helper()
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			buf := make([]byte, 1024)
			n, err := r.Read(buf)
			if err != nil {
				tb.Fatalf("wait %q, but received %q: %v", s, text, err)
			}
			curr := string(buf[:n])
			text += curr
			tb.Logf("reply: %q", curr)
			if strings.LastIndex(text, s) >= 0 {
				return
			}
		}
	}
}
