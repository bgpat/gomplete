package e2e

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const timeout = 5 * time.Second

func writeAndWait(ctx context.Context, r io.ReadWriter, s string) error {
	if _, err := io.WriteString(r, s); err != nil {
		return errors.WithStack(err)
	}
	_, err := waitString(ctx, r, s)
	return errors.WithStack(err)
}

func waitString(ctx context.Context, r io.Reader, s string) (text string, err error) {
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
			text += string(buf[:n])
			if strings.LastIndex(text, s) >= 0 {
				return text, nil
			}
		}
	}
}
