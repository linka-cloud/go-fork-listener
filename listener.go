// Copyright 2022 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fork

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	defaultEnv = "GO_FORK_CHILD"
)

var (
	_ net.Listener = (*forkedListener)(nil)

	ErrClosed = errors.New("connection closed")
)

func IsClosed(err error) bool {
	return errors.Is(err, ErrClosed)
}

type Listener interface {
	net.Listener
	Start() error
	Wait() error
	IsChild() bool
	IsParent() bool
}

// Listen creates a new listener that will fork a child process when a new connection is accepted.
func Listen(network, address string, opts ...Option) (Listener, error) {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	if os.Getenv(o.env) == "1" {
		return NewChildListener()
	}
	inner, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	switch inner.(type) {
	case *net.TCPListener, *net.UnixListener:
	default:
		return nil, fmt.Errorf("unsupported listener type %T", inner)
	}
	return &listener{inner: inner, errs: make(chan error, 1), opts: &o}, nil
}
