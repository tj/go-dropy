// Package dropy implements a higher-level Dropbox API on top of go-dropbox.
package dropy

import (
	"io"
	"io/ioutil"
	"os"
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

// Client wraps dropbox.Client to provide higher level sugar.
type Client struct {
	*dropbox.Client
}

// New client.
func New(d *dropbox.Client) *Client {
	return &Client{
		Client: d,
	}
}

// Stat returns file and directory meta-data for `name`.
func (c *Client) Stat(name string) (os.FileInfo, error) {
	out, err := c.Files.GetMetadata(&dropbox.GetMetadataInput{
		Path: name,
	})

	if err != nil {
		return nil, err
	}

	return &FileInfo{&out.Metadata}, nil
}

// List returns entries in dir `name`. Up to `n` entries, or all when `n` <= 0.
func (c *Client) List(name string, n int) (ents []os.FileInfo, err error) {
	var cursor string

	if n <= 0 {
		n = -1
	}

	for {
		var out *dropbox.ListFolderOutput

		if cursor == "" {
			out, err = c.Files.ListFolder(&dropbox.ListFolderInput{Path: name})
			cursor = out.Cursor
		} else {
			out, err = c.Files.ListFolderContinue(&dropbox.ListFolderContinueInput{cursor})
			cursor = out.Cursor
		}

		if err != nil {
			return
		}

		for _, ent := range out.Entries {
			ents = append(ents, &FileInfo{ent})
		}

		if n >= 0 && len(ents) >= n {
			ents = ents[:n]
			break
		}

		if !out.HasMore {
			break
		}
	}

	if n >= 0 && len(ents) == 0 {
		err = io.EOF
		return
	}

	return
}

// ListAll returns all entries in dir `name`.
func (c *Client) ListAll(name string) ([]os.FileInfo, error) {
	return c.List(name, 0)
}

// Open returns a File for reading and writing.
func (c *Client) Open(name string) *File {
	return &File{
		Name: name,
		c:    c,
	}
}

// ReadAll returns the contents of `name`.
func (c *Client) ReadAll(name string) ([]byte, error) {
	f := c.Open(name)
	defer f.Close()
	return ioutil.ReadAll(f)
}
