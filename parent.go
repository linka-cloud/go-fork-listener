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
	"net"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// FileConn is a connection that can be transferred to a child process using file descriptors.
type FileConn interface {
	net.Conn
	File() (*os.File, error)
}

type listener struct {
	inner net.Listener
	wg    sync.WaitGroup
	errs  chan error

	opts *options
}

// Start starts the listener in the background.
func (l *listener) Start() error {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.Accept()
	}()
	return nil
}

// Wait waits for the listener to exit
func (l *listener) Wait() error {
	l.wg.Wait()
	select {
	case err := <-l.errs:
		return err
	default:
		return nil
	}
}

// Run starts the listener in the background and waits for it to exit.
func (l *listener) Run() error {
	if err := l.Start(); err != nil {
		return err
	}
	return l.Wait()
}

// IsChild returns true if the current process is a forked child process.
func (l *listener) IsChild() bool {
	return false
}

// IsParent returns true if the current process is the main (parent) process.
func (l *listener) IsParent() bool {
	return true
}

// Accept implements the net.Listener interface.
// It will wait for the underlying listener to accept a new connection
// then transfer it to the child process.
// It actually never returns a connection.
func (l *listener) Accept() (net.Conn, error) {
	for {
		conn, err := l.inner.Accept()
		if err != nil {
			l.errs <- err
			return nil, err
		}
		go func() {
			defer conn.Close()
			var fc FileConn
			switch c := conn.(type) {
			case *net.TCPConn:
				fc = c
			case *net.UnixConn:
				fc = c
			default:
				logrus.Fatalf("unsupported connection type %T", c)
			}
			if err := l.fork(fc); err != nil {
				logrus.Errorf("fork failed: %v", err)
			}
		}()
	}
}

// Addr implements the net.Listener interface.
func (l *listener) Addr() net.Addr {
	return l.inner.Addr()
}

// Close implements the net.Listener interface.
func (l *listener) Close() error {
	return l.inner.Close()
}
