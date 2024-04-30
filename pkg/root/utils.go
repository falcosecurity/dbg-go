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

package root

import (
	"fmt"
	"path/filepath"
)

func (f *FsLooper) LoopFiltered(opts Options, message, tag string, worker RowWorker) error {
	configNameGlob := opts.Target.toGlob()
	for _, driverVersion := range opts.DriverVersion {
		path := f.builder(opts, driverVersion, configNameGlob)
		files, err := filepath.Glob(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			Printer.Logger.Info(message,
				Printer.Logger.Args(tag, file))
			if opts.DryRun {
				Printer.Logger.Info("skipping because of dry-run.")
				return nil
			}
			err = worker(driverVersion, file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func BuildConfigPath(opts Options, driverVersion, configName string) string {
	return fmt.Sprintf(configPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture.ToNonDeb(),
		configName)
}

func BuildOutputPath(opts Options, driverVersion, outputName string) string {
	fullName := ""
	if outputName != "" {
		// only add "drivername_" prefix when outputName is not empty,
		// ie: when we are not generating a folder path.
		fullName = opts.DriverName + "_" + outputName
	}

	return fmt.Sprintf(outputPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture.ToNonDeb(),
		fullName)
}
