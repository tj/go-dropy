package dropy

import (
	"io"

	"github.com/tj/go-dropbox"
)

// File implements an io.ReadWriteCloser for Dropbox files.
type File struct {
	Name string
	c    *Client
	r    io.ReadCloser
}

// Read implementation, note that the first call to this method
// triggers the download, seeking is currently not supported.
func (f *File) Read(b []byte) (int, error) {
	if f.r == nil {
		if err := f.download(); err != nil {
			return 0, err
		}
	}

	return f.r.Read(b)
}

// download the file.
func (f *File) download() error {
	out, err := f.c.Files.Download(&dropbox.DownloadInput{f.Name})
	if err != nil {
		return nil
	}

	f.r = out.Body
	return nil
}

// Close the file.
func (f *File) Close() error {
	return nil
}
