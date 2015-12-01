package dropy

import (
	"io/ioutil"
	"os"
	"strings"
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

func TestClient_List(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.List("/list")
	assert.NoError(t, err)
	assert.Equal(t, 5000, len(ents))
}

func TestClient_ListN_missing(t *testing.T) {
	t.Parallel()
	c := client()
	_, err := c.ListN("/notfound", 0)
	assert.Error(t, err)
}

func TestClient_ListN_zero(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListN("/list", 0)
	assert.NoError(t, err)
	assert.Equal(t, 5000, len(ents))
}

func TestClient_ListN_subzero(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListN("/list", -5)
	assert.NoError(t, err)
	assert.Equal(t, 5000, len(ents))
}

func TestClient_ListN_count(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListN("/list", 1234)
	assert.NoError(t, err)
	assert.Equal(t, 1234, len(ents))
}

func TestClient_ListFilter(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListFilter("/list-types", func(info os.FileInfo) bool {
		return info.IsDir()
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ents))
}

func TestClient_ListFolders(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListFolders("/list-types")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ents))
	assert.Equal(t, "one", ents[0].Name())
}

func TestClient_ListFiles(t *testing.T) {
	t.Parallel()
	c := client()
	ents, err := c.ListFiles("/list-types")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ents))
	assert.Equal(t, "one.txt", ents[0].Name())
}

func TestClient_Open(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")

	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, "world", string(b))
}

func TestCient_Open_missing(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/dev/null")

	_, err := ioutil.ReadAll(f)
	assert.EqualError(t, err, "open /dev/null: no such file or directory")
}

func TestClient_Read(t *testing.T) {
	t.Parallel()
	c := client()
	b, err := c.Read("/hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "world", string(b))
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/delete.txt")
	f.Write([]byte("Hello World"))
	assert.NoError(t, f.Close())

	assert.NoError(t, c.Delete("/delete.txt"))
}

func TestClient_Search(t *testing.T) {
	t.Parallel()
	c := client()

	list, err := c.Search("/list", "100")
	assert.NoError(t, err)

	assert.Equal(t, 11, len(list))
}

func TestClient_Search_more(t *testing.T) {
	t.Parallel()
	c := client()

	list, err := c.Search("/list", "10")
	assert.NoError(t, err)

	assert.Equal(t, 111, len(list))
}

func TestClient_Upload(t *testing.T) {
	t.Parallel()

	c := client()
	err := c.Upload("/upload-1.txt", strings.NewReader("one"))
	assert.NoError(t, err, "error uploading")

	b, err := c.Read("/upload-1.txt")
	assert.NoError(t, err, "error reading")

	assert.Equal(t, "one", string(b))
}
