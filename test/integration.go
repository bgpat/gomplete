// +build integration

package test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

const timeout = 5 * time.Second

// WriteAndWait writes the string and waits the terminal response.
func WriteAndWait(ctx context.Context, tb *testing.TB, r io.ReadWriter, s string) error {
	tb.Helper()
	if _, err := io.WriteString(r, s); err != nil {
		return errors.WithStack(err)
	}
	_, err := waitString(ctx, tb, r, s)
	return errors.WithStack(err)
}

// WaitString waits the terminal still the presented string is printed.
func WaitString(ctx context.Context, tb *testing.TB, r io.Reader, s string) (text string, err error) {
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
				return text, errors.Wrapf(err, "wait %q, but received %q", s, text)
			}
			curr := string(buf[:n])
			text += curr
			tb.Log(curr)
			if strings.LastIndex(text, s) >= 0 {
				return text, nil
			}
		}
	}
}
