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

package cleanup

import (
	"os"

	"github.com/falcosecurity/dbg-go/pkg/root"
)

type fileCleaner struct {
	root.Looper
}

func NewFileCleaner() Cleaner {
	return &fileCleaner{Looper: root.NewFsLooper(root.BuildConfigPath)}
}

func (f *fileCleaner) Info() string {
	return "cleaning up local config files"
}

func (f *fileCleaner) Cleanup(opts Options) error {
	return f.LoopFiltered(opts.Options, "removing file", "config", func(driverVersion, configPath string) error {
		return os.Remove(configPath)
	})
}
