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

	assert.Equal(t, "world", string(b))
}

func TestFile_Close(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")
	assert.NoError(t, f.Close())
}

func TestFile_Close_inval(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")
	assert.NoError(t, f.Close())
	assert.EqualError(t, f.Close(), "close /hello.txt: invalid argument")
}

func TestFile_Read(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")

	b := make([]byte, 5)
	n, err := f.Read(b)
	assert.Equal(t, 5, n)
	assert.EqualError(t, err, "EOF")
	assert.Equal(t, "world", string(b))

	assert.NoError(t, f.Close())
}

func TestFile_Write(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello-world-1.txt")

	_, err := f.Write([]byte("Hello"))
	assert.NoError(t, err)

	_, err = f.Write([]byte(" Wor"))
	assert.NoError(t, err)

	_, err = f.Write([]byte("ld"))
	assert.NoError(t, err)

	assert.NoError(t, f.Close())

	b, err := c.Read("/hello-world-1.txt")
	assert.NoError(t, err)

	assert.Equal(t, "Hello World", string(b))
}
