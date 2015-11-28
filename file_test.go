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

func TestFile_Close(t *testing.T) {
	t.Parallel()
	c := client()

	f := c.Open("/hello.txt")
	assert.NoError(t, f.Close())
}
