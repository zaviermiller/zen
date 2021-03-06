package vm

import (
	"io"

	u "github.com/zaviermiller/zen/internal/util"
)

// Custom updater for custom read func
type zenUpdateWriter struct {
	io.Reader
	total int64
	size  int
}

// Custom read func which implements loader
func (z *zenUpdateWriter) Read(p []byte) (int, error) {
	n, err := z.Reader.Read(p)
	z.total += int64(n)

	if err == nil {
		u.PrintSimpleLoader(int(z.total), z.size, "Downloading update...")
	}

	return n, err
}
