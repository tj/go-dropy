package dropy_test

import (
	"io"
	"os"
	"strings"

	"github.com/segmentio/go-env"
	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
)

// Upload and read a file.
func Example() {
	token := env.MustGet("DROPBOX_ACCESS_TOKEN")
	client := dropy.New(dropbox.New(dropbox.NewConfig(token)))

	file := client.Open("/demo.txt")
	io.Copy(file, strings.NewReader("Hello World"))

	io.Copy(os.Stdout, file)
	// Output: Hello World

	file.Close()
}
