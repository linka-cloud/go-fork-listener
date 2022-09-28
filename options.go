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
	"io"
	"os"
)

type Option func(o *options)

func WithChildEnvName(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

func WithEnv(env []string) Option {
	return func(o *options) {
		o.extraEnv = env
	}
}

func WithStdin(r io.Reader) Option {
	return func(o *options) {
		o.stdin = r
	}
}

func WithStdout(w io.Writer) Option {
	return func(o *options) {
		o.stdout = w
	}
}

func WithStderr(w io.Writer) Option {
	return func(o *options) {
		o.stderr = w
	}
}

func WithUID(uid int) Option {
	return func(o *options) {
		o.uid = uid
	}
}

func WithGID(gid int) Option {
	return func(o *options) {
		o.gid = gid
	}
}

type options struct {
	env      string
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
	extraEnv []string
	uid      int
	gid      int
}

var defaultOptions = options{
	env:    defaultEnv,
	stdin:  os.Stdin,
	stdout: os.Stdout,
	stderr: os.Stderr,
	uid:    -1,
	gid:    -1,
}
