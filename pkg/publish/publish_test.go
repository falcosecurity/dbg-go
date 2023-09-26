package publish

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/dbg-go/pkg/root"
	s3utils "github.com/falcosecurity/dbg-go/pkg/utils/s3"
	testutils "github.com/falcosecurity/dbg-go/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPublish(t *testing.T) {
	// This client will be used by the Run action
	testClient = testutils.S3CreateTestBucket(t, nil)

	outputPath := root.BuildOutputPath(root.Options{
		RepoRoot:     "./test",
		Architecture: "amd64",
	}, "5.0.1+driver", "")
	err := os.MkdirAll(outputPath, 0700)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll("./test")
	})

	// Create a fake kernel module object to be uploaded to test bucket
	d1 := []byte("TEST\n")
	err = os.WriteFile(outputPath+"/falco_almalinux_4.18.0-425.10.1.el8_7.x86_64_1.ko", d1, 0644)
	assert.NoError(t, err)

	// Fetch an existing object metadata
	realClient, err := s3utils.NewClient(true)
	assert.NoError(t, err)
	object, err := realClient.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s3utils.S3Bucket),
		Key:    aws.String("driver/5.0.1+driver/x86_64/falco_almalinux_4.18.0-425.10.1.el8_7.x86_64_1.ko"),
	})
	assert.NoError(t, err)

	// Run our action to upload our driver object
	err = Run(Options{
		Options: root.Options{
			RepoRoot:      "./test",
			Architecture:  "amd64",
			DriverName:    "falco",
			DriverVersion: []string{"5.0.1+driver"},
			Target: root.Target{
				Distro:        "almalinux",
				KernelRelease: "4.18.0-425*",
				KernelVersion: "1",
			},
		},
	})
	assert.NoError(t, err)

	// Fetch test object metadata
	testObject, err := testClient.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s3utils.S3Bucket),
		Key:    aws.String("driver/5.0.1+driver/x86_64/falco_almalinux_4.18.0-425.10.1.el8_7.x86_64_1.ko"),
	})
	assert.NoError(t, err)

	// Check that published object has correct metadata
	assert.Equal(t, testObject.ServerSideEncryption, object.ServerSideEncryption)
	assert.Equal(t, testObject.AcceptRanges, object.AcceptRanges)
	assert.Equal(t, testObject.ObjectLockMode, object.ObjectLockMode)
	assert.Equal(t, testObject.ArchiveStatus, object.ArchiveStatus)
	assert.Equal(t, testObject.BucketKeyEnabled, object.BucketKeyEnabled)
	assert.Equal(t, testObject.CacheControl, object.CacheControl)
	assert.Equal(t, testObject.ContentDisposition, object.ContentDisposition)
	assert.Equal(t, testObject.ContentEncoding, object.ContentEncoding)
	assert.Equal(t, testObject.ContentType, object.ContentType)
	assert.Equal(t, testObject.Expiration, object.Expiration)
	assert.Equal(t, testObject.Expires, object.Expires)
	assert.Equal(t, testObject.Metadata, object.Metadata)
	assert.Equal(t, testObject.DeleteMarker, object.DeleteMarker)
	assert.Equal(t, testObject.MissingMeta, object.MissingMeta)
	assert.Equal(t, testObject.ObjectLockLegalHoldStatus, object.ObjectLockLegalHoldStatus)
	assert.Equal(t, testObject.ObjectLockRetainUntilDate, object.ObjectLockRetainUntilDate)
	assert.Equal(t, testObject.PartsCount, object.PartsCount)
	assert.Equal(t, testObject.ReplicationStatus, object.ReplicationStatus)
	assert.Equal(t, testObject.RequestCharged, object.RequestCharged)
	assert.Equal(t, testObject.Restore, object.Restore)
	assert.Equal(t, testObject.StorageClass, object.StorageClass)
}
