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
		if err := lis.Start(); err != nil {
			logrus.Fatal(err)
		}
		if err := lis.Wait(); err != nil {
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
