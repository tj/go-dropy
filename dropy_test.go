package dropy

import (
	"testing"

	"github.com/segmentio/go-env"
	"github.com/stretchr/testify/assert"
	"github.com/tj/go-dropbox"
)

func client() *Client {
	token := env.MustGet("DROPBOX_ACCESS_TOKEN")
	return New(dropbox.New(dropbox.NewConfig(token)))
}

func TestClient_Stat(t *testing.T) {
	t.Parallel()
	c := client()
	info, err := c.Stat("/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, false, info.IsDir())
	assert.Equal(t, false, info.Mode().IsDir())
	assert.Equal(t, true, info.Mode().IsRegular())
	assert.Equal(t, "hello.txt", info.Name())
	assert.Equal(t, int64(5), info.Size())
}

func TestClient_ListAll(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListAll("/list")
	assert.NoError(t, err)
	assert.Equal(t, 5000, len(ents))
}

func TestClient_List_zero(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.List("/list", 0)
	assert.NoError(t, err)
	assert.Equal(t, 5000, len(ents))
}

func TestClient_List_subzero(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.List("/list", -5)
	assert.NoError(t, err)
	assert.Equal(t, 5000, len(ents))
}

func TestClient_List_count(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.List("/list", 1234)
	assert.NoError(t, err)
	assert.Equal(t, 1234, len(ents))
}

func TestClient_ReadAll(t *testing.T) {
	t.Parallel()
	c := client()
	b, err := c.ReadAll("/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "whoop", string(b))
}
