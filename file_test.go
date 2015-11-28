package dropy

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Open(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")

	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, "whoop", string(b))
}

func TestFile_Sync(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello-world.txt")

	_, err := f.Write([]byte("Hello"))
	assert.NoError(t, err)

	_, err = f.Write([]byte(" World"))
	assert.NoError(t, err)

	assert.NoError(t, f.Sync())
	assert.NoError(t, f.Close())

	b, err := c.ReadAll("/hello-world.txt")
	assert.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}

func TestFile_Close_read(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")
	assert.NoError(t, f.Close())
}

func TestFile_Close_write(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello-world.txt")

	_, err := f.Write([]byte("Hello"))
	assert.NoError(t, err)

	_, err = f.Write([]byte(" World"))
	assert.NoError(t, err)

	assert.NoError(t, f.Close())

	b, err := c.ReadAll("/hello-world.txt")
	assert.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}
