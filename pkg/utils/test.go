//go:build test_all

package utils

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TODO: cleanup_s3 tests

func RunTestParsingLogs(t *testing.T, runTest func() error, parsedMsg interface{}, parsingCB func() bool) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(io.Writer(&buf), nil))
	slog.SetDefault(logger)

	err := runTest()
	assert.NoError(t, err)

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), parsedMsg)
		assert.NoError(t, err)
		if parsingCB() == false {
			break
		}
	}
}

func PreCreateFolders(repoRoot, architecture string, driverVersionsToBeCreated []string) error {
	for _, driverVersion := range driverVersionsToBeCreated {
		configPath := fmt.Sprintf(root.ConfigPathFmt,
			repoRoot,
			driverVersion,
			architecture,
			"")
		err := os.MkdirAll(configPath, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

// SliceDifference returns the elements in `a` that aren't in `b`.
func SliceDifference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func S3CreateTestBucket(t *testing.T) *s3.Client {
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	t.Cleanup(func() {
		ts.Close()
	})

	// Difference in configuring the client

	// Setup a new config
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("KEY", "SECRET", "SESSION")),
		config.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: ts.URL}, nil
			}),
		),
	)

	// Create an Amazon S3 v2 client, important to use o.UsePathStyle
	// alternatively change local DNS settings, e.g., in /etc/hosts
	// to support requests to http://<bucketname>.127.0.0.1:32947/...
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Create bucket
	_, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(S3Bucket),
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, _ = client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
			Bucket: aws.String(S3Bucket),
		})
	})

	// Create some test keys
	keysToBeCreated := []string{
		"driver/1.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
		"driver/1.0.0+driver/x86_64/falco_amazonlinux2022_5.10.96-90.460.amzn2022.x86_64_1.o",
		"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.o",
		"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.ko",
		"driver/2.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
		"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
		"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
	}
	for _, key := range keysToBeCreated {
		_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(S3Bucket),
			Key:    aws.String(key),
		})
		assert.NoError(t, err)
	}
	return client
}
