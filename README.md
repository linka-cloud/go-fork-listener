# go-fork-listener

[![Language: Go](https://img.shields.io/badge/lang-Go-6ad7e5.svg?style=flat-square&logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/go-fork-listener.svg)](https://pkg.go.dev/go.linka.cloud/go-fork-listener)

Go package for listening on a socket and forking a detached child process on each connection.

This is useful if you need to be able to restart a process without interrupting any connections, like *sshd*.

**Only unix like systems are supported as the connection is passed to the child process via file descriptor.**

## Usage

```go
import fork "go.linka.cloud/go-fork-listener"
```

## Example


### Simple http server

```go
package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"go.linka.cloud/go-fork-listener"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Infof("running with pid %d", os.Getpid())
	lis, err := fork.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatal(err)
	}
	defer lis.Close()
	if err := http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello from child process!\n")); err != nil {
			logrus.Error(err)
		}
	})); err != nil && !fork.IsClosed(err) {
		logrus.Fatal(err)
	}
}
```

### Fork aware http server

```go
package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"go.linka.cloud/go-fork-listener"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Infof("running with pid %d", os.Getpid())
	lis, err := fork.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatal(err)
	}
	defer lis.Close()

	if lis.IsParent() {
		// do some configuration checks or other things that need to be done only once
		if err := lis.Run(); err != nil {
			logrus.Fatal(err)
		}
		return
	}

	// do some handler dependencies initialization

	if err := http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello from child process!\n")); err != nil {
			logrus.Error(err)
		}
	})); err != nil && !fork.IsClosed(err) {
		logrus.Fatal(err)
	}
}

```
