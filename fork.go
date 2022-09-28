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
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func (l *listener) fork(conn FileConn) error {
	logrus.Debugf("%s: spawning child process", conn.RemoteAddr())
	uid, gid := uint32(os.Getuid()), uint32(os.Getgid())
	if l.opts.uid != -1 {
		uid = uint32(l.opts.uid)
	}
	if l.opts.gid != -1 {
		gid = uint32(l.opts.gid)
	}
	f, err := conn.File()
	if err != nil {
		return err
	}
	defer f.Close()
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:     true,
		Credential: &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true},
	}
	cmd.ExtraFiles = []*os.File{f}
	cmd.Stdin = l.opts.stdin
	cmd.Stdout = l.opts.stdout
	cmd.Stderr = l.opts.stderr
	env := append(os.Environ(), l.opts.extraEnv...)
	cmd.Env = append(env, fmt.Sprintf("%s=1", l.opts.env))
	if err := cmd.Start(); err != nil {
		return err
	}
	defer logrus.Debugf("%s: child process released", conn.RemoteAddr())
	if err := cmd.Wait(); err != nil {
		logrus.Errorf("%s: child process exited with error: %v", conn.RemoteAddr(), err)
		return err
	}
	return nil
}
