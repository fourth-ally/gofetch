package infrastructure

import (
	"io"

	"github.com/nkechr1s/gofetch/domain/contracts"
)

// progressReader wraps an io.Reader to track progress.
type progressReader struct {
	reader      io.Reader
	total       int64
	transferred int64
	callback    contracts.ProgressCallback
}

// Read implements io.Reader interface with progress tracking.
func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.transferred += int64(n)

	if pr.callback != nil {
		pr.callback(pr.transferred, pr.total)
	}

	return n, err
}
