// Package dropy implements a higher-level Dropbox API on top of go-dropbox.
package dropy

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/tj/go-dropbox"
)

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

// ListN returns entries in dir `name`. Up to `n` entries, or all when `n` <= 0.
func (c *Client) ListN(name string, n int) (ents []os.FileInfo, err error) {
	var cursor string

	if n <= 0 {
		n = -1
	}

	for {
		var out *dropbox.ListFolderOutput

		if cursor == "" {
			if out, err = c.Files.ListFolder(&dropbox.ListFolderInput{Path: name}); err != nil {
				return
			}
			cursor = out.Cursor
		} else {
			if out, err = c.Files.ListFolderContinue(&dropbox.ListFolderContinueInput{cursor}); err != nil {
				return
			}
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

// List returns all entries in dir `name`.
func (c *Client) List(name string) ([]os.FileInfo, error) {
	return c.ListN(name, 0)
}

// ListFilter returns all entries in dir `name` filtered by `filter`.
func (c *Client) ListFilter(name string, filter func(info os.FileInfo) bool) (ret []os.FileInfo, err error) {
	ents, err := c.ListN(name, 0)
	if err != nil {
		return
	}

	for _, ent := range ents {
		if filter(ent) {
			ret = append(ret, ent)
		}
	}

	return
}

// ListFolders returns all folders in dir `name`.
func (c *Client) ListFolders(name string) ([]os.FileInfo, error) {
	return c.ListFilter(name, func(info os.FileInfo) bool {
		return info.IsDir()
	})
}

// Open returns a File for reading and writing.
func (c *Client) Open(name string) *File {
	r, w := io.Pipe()
	return &File{
		Name:  name,
		c:     c,
		pipeR: r,
		pipeW: w,
	}
}

// ReadAll returns the contents of `name`.
func (c *Client) ReadAll(name string) ([]byte, error) {
	f := c.Open(name)
	defer f.Close()
	return ioutil.ReadAll(f)
}
