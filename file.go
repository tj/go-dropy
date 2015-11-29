package dropy

import (
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/tj/go-dropbox"
)

// FileInfo wraps Dropbox file MetaData to implement os.FileInfo.
type FileInfo struct {
	meta *dropbox.Metadata
}

// Name of the file.
func (f *FileInfo) Name() string {
	return f.meta.Name
}

// Size of the file.
func (f *FileInfo) Size() int64 {
	return int64(f.meta.Size)
}

// IsDir returns true if the file is a directory.
func (f *FileInfo) IsDir() bool {
	return f.meta.Tag == "folder"
}

// Sys is not implemented.
func (f *FileInfo) Sys() interface{} {
	return nil
}

// ModTime returns the modification time.
func (f *FileInfo) ModTime() time.Time {
	return f.meta.ServerModified
}

// Mode returns the file mode flags.
func (f *FileInfo) Mode() os.FileMode {
	var m os.FileMode

	if f.IsDir() {
		m |= os.ModeDir
	}

	return m
}

// File implements an io.ReadWriteCloser for Dropbox files.
type File struct {
	Name    string
	closed  bool
	writing bool
	reader  io.ReadCloser
	pipeR   *io.PipeReader
	pipeW   *io.PipeWriter
	c       *Client
}

// Read implements io.Reader
func (f *File) Read(b []byte) (int, error) {
	if f.reader == nil {
		if err := f.download(); err != nil {
			return 0, err
		}
	}

	return f.reader.Read(b)
}

// Write implements io.Writer.
func (f *File) Write(b []byte) (int, error) {
	if !f.writing {
		f.writing = true

		go func() {
			_, err := f.c.Files.Upload(&dropbox.UploadInput{
				Mode:   dropbox.WriteModeOverwrite,
				Path:   f.Name,
				Mute:   true,
				Reader: f.pipeR,
			})

			f.pipeR.CloseWithError(err)
		}()
	}

	return f.pipeW.Write(b)
}

// Close implements io.Closer.
func (f *File) Close() error {
	if f.closed {
		return &os.PathError{"close", f.Name, syscall.EINVAL}
	}
	f.closed = true

	if f.writing {
		if err := f.pipeW.Close(); err != nil {
			return err
		}
	}

	if f.reader != nil {
		if err := f.reader.Close(); err != nil {
			return err
		}
	}

	return nil
}

// download the file.
func (f *File) download() error {
	out, err := f.c.Files.Download(&dropbox.DownloadInput{f.Name})
	if err != nil {
		if strings.HasPrefix(err.Error(), "path/not_found/") {
			return &os.PathError{"open", f.Name, syscall.ENOENT}
		}
		return err
	}

	f.reader = out.Body
	return nil
}
