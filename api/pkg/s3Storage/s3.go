package s3Storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewStorageS3(host, region, accessKey, secretKey string) (*s3.S3, error) {
	s3Session := s3.New(session.Must(session.NewSession(&aws.Config{
		Endpoint:    &host,
		Region:      &region,
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})))

	return s3Session, nil
}
