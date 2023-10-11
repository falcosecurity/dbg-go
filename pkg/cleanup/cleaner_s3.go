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
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3utils "github.com/falcosecurity/dbg-go/pkg/utils/s3"
)

type s3Cleaner struct {
	*s3utils.Client
}

func NewS3Cleaner() (Cleaner, error) {
	client, err := s3utils.NewClient(false)
	if err != nil {
		return nil, err
	}
	return &s3Cleaner{Client: client}, nil
}

func (s *s3Cleaner) Info() string {
	return "cleaning up remote driver files"
}

func (s *s3Cleaner) Cleanup(opts Options) error {
	return s.LoopFiltered(opts.Options, "cleaning up remote driver file", "key", func(driverVersion, key string) error {
		_, err := s.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(s3utils.S3Bucket),
			Key:    aws.String(key),
		})
		return err
	})
}
