package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/godev111222333/capstone-backend/src/misc"
)

type CredentialsProvider struct {
	AccessKey       string
	SecretAccessKey string
}

func NewCredentialsProvider(accessKey, secretAccessKey string) *CredentialsProvider {
	return &CredentialsProvider{accessKey, secretAccessKey}
}

func (c *CredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.AccessKey,
		SecretAccessKey: c.SecretAccessKey,
	}, nil
}

type S3Store struct {
	Config *misc.AWSConfig
	Client *s3.Client
}

func NewS3Store(cfg *misc.AWSConfig) *S3Store {
	client := s3.New(s3.Options{
		Credentials: NewCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey),
		Region:      cfg.Region,
	})
	return &S3Store{Client: client, Config: cfg}
}
