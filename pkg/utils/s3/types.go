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
