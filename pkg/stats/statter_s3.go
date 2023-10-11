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

package stats

import (
	"log/slog"
	"strings"

	"github.com/falcosecurity/dbg-go/pkg/root"
	s3utils "github.com/falcosecurity/dbg-go/pkg/utils/s3"
)

type s3Statter struct {
	*s3utils.Client
}

func NewS3Statter() (Statter, error) {
	client, err := s3utils.NewClient(true)
	if err != nil {
		return nil, err
	}
	return &s3Statter{Client: client}, nil
}

func (f *s3Statter) Info() string {
	return "gathering stats for remote drivers"
}

func (s *s3Statter) GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error) {
	slog.SetDefault(slog.With("bucket", s3utils.S3Bucket))

	driverStatsByVersion := make(driverStatsByDriverVersion)
	err := s.LoopFiltered(opts, "computing stats", "key", func(driverVersion, key string) error {
		dStats := driverStatsByVersion[driverVersion]
		if strings.HasSuffix(key, ".ko") {
			dStats.NumModules++
		} else if strings.HasSuffix(key, ".o") {
			dStats.NumProbes++
		}
		driverStatsByVersion[driverVersion] = dStats
		return nil
	})
	return driverStatsByVersion, err
}
