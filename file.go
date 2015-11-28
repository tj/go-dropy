package dropy

import (
	"bytes"
	"io"

	"github.com/tj/go-dropbox"
)

// File implements an io.ReadWriteCloser for Dropbox files.
type File struct {
	Name string
	c    *Client
	w    bytes.Buffer
	r    io.ReadCloser
}

// Read implements io.Reader
//
// Note that the first call to this method triggers
// the download, seeking is currently not supported.
func (f *File) Read(b []byte) (int, error) {
	if f.r == nil {
		if err := f.download(); err != nil {
			return 0, err
		}
	}

	return f.r.Read(b)
}

// Write implements io.Writer.
//
// Note that the upload occurs when the Close
// or Sync methods are invoked, until then
// the contents are buffered in-memory.
func (f *File) Write(b []byte) (int, error) {
	return f.w.Write(b)
}

// Sync the file to Dropbox.
func (f *File) Sync() error {
	_, err := f.c.Files.Upload(&dropbox.UploadInput{
		Mode:   dropbox.WriteModeOverwrite,
		Path:   f.Name,
		Mute:   true,
		Reader: bytes.NewBuffer(f.w.Bytes()),
	})

	f.w.Reset()

	return err
}

// Close implements io.Closer.
func (f *File) Close() error {
	if f.w.Len() > 0 {
		return f.Sync()
	}

	return nil
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
