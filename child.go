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
	"net"
	"os"
	"sync"
)

func NewChildListener() (Listener, error) {
	f := os.NewFile(3, "conn")
	if f == nil {
		return nil, errors.New("invalid connection file descriptor")
	}
	c, err := net.FileConn(f)
	if err != nil {
		return nil, err
	}
	ch := make(chan net.Conn, 1)
	lis := &forkedListener{
		ch:     ch,
		addr:   c.LocalAddr(),
		closed: make(chan struct{}),
	}
	conn := &conn{Conn: c, lis: lis}
	ch <- conn
	return lis, nil
}

type forkedListener struct {
	o      sync.Once
	ch     chan net.Conn
	addr   net.Addr
	closed chan struct{}
}

func (l *forkedListener) Start() error {
	return errors.New("cannot start a child listener")
}

func (l *forkedListener) Wait() error {
	return errors.New("cannot wait on a child listener")
}

func (l *forkedListener) Run() error {
	return errors.New("cannot run a child listener")
}

func (l *forkedListener) IsChild() bool {
	return true
}

func (l *forkedListener) IsParent() bool {
	return false
}

func (l *forkedListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.ch:
		return conn, nil
	case <-l.closed:
		return nil, ErrClosed
	}
}

func (l *forkedListener) Close() error {
	l.o.Do(func() {
		close(l.closed)
	})
	return nil
}

func (l *forkedListener) Addr() net.Addr {
	return l.addr
}

type conn struct {
	net.Conn
	lis *forkedListener
}

func (c *conn) Close() error {
	err := c.Conn.Close()
	c.lis.o.Do(func() {
		close(c.lis.closed)
	})
	return err
}
