
[![GoDoc](https://godoc.org/github.com/tj/go-dropy?status.svg)](https://godoc.org/github.com/tj/go-dropy) [![Build Status](https://semaphoreci.com/api/v1/projects/486c0583-68ae-465b-a25e-422dd8760f6e/617435/badge.svg)](https://semaphoreci.com/tj/go-dropy)


# Dropy

 High level Dropbox v2 client for Go built on top of [go-dropbox](https://github.com/tj/go-dropbox).

## Example

```go
token := os.Getenv("DROPBOX_ACCESS_TOKEN")
client := dropy.New(dropbox.New(dropbox.NewConfig(token)))

client.Upload("/demo.txt", strings.NewReader("Hello World"))
io.Copy(os.Stdout, client.Open("/demo.txt"))
```

## Testing

 To manually run tests use the test account access token:

```
$ export DROPBOX_ACCESS_TOKEN=oENFkq_oIVAAAAAAAAAABqI2Nor2e9_ORA3oAZDQexMgJocCQX4aOFXZuDc1t-Sx
$ go test -v
```

# License

 MIT