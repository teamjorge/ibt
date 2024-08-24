package headers

import "io"

// Reader provides all IO functions necessary for parsing ibt files.
type Reader interface {
	io.Reader
	io.ReaderAt
	io.Closer
	io.ReadSeeker
}
