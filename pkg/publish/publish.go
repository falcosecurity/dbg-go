// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package publish

import (
	"github.com/falcosecurity/dbg-go/pkg/root"
	s3utils "github.com/falcosecurity/dbg-go/pkg/utils/s3"
)

// Used by tests
var testClient *s3utils.Client

func Run(opts Options) error {
	root.Printer.Logger.Info("publishing drivers")
	var (
		client *s3utils.Client
		err    error
	)
	if testClient == nil {
		client, err = s3utils.NewClient(false)
		if err != nil {
			return err
		}
	} else {
		client = testClient
	}
	looper := root.NewFsLooper(root.BuildOutputPath)
	return looper.LoopFiltered(opts.Options, "publishing", "driver", func(driverVersion, path string) error {
		return client.PutDriver(opts.Options, driverVersion, path)
	})
}
