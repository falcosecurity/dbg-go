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

package s3utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	*s3.Client
}

func NewClient(readOnly bool) (*Client, error) {
	var (
		cfg aws.Config
		err error
	)
	if !readOnly {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(s3Region))
		if err != nil {
			return nil, err
		}
	} else {
		cfg = aws.Config{
			Region:      s3Region,
			Credentials: aws.AnonymousCredentials{},
		}
	}
	return &Client{Client: s3.NewFromConfig(cfg)}, nil
}
