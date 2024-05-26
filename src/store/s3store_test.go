package store

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestS3Store_AllOperations(t *testing.T) {
	t.Parallel()

	t.Run("list all objects", func(t *testing.T) {
		t.Parallel()

		key := uuid.NewString()
		pwd, _ := os.Getwd()
		fileReader, err := os.ReadFile(filepath.Join(pwd, "/../../etc/ava.png"))
		require.NoError(t, err)
		_, err = TestS3Store.Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: &TestS3Store.Config.Bucket,
			Key:    aws.String(key + ".png"),
			Body:   bytes.NewReader(fileReader),
			ACL:    types.ObjectCannedACLPublicRead,
		})
		require.NoError(t, err)
	})
}
