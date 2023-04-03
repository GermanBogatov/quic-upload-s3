package storage

import "github.com/aws/aws-sdk-go/service/s3"

type Storage struct {
	db     *s3.S3
	bucket string
}

func NewStorage(db *s3.S3, bucket string) IStorage {
	return &Storage{
		db:     db,
		bucket: bucket,
	}
}

func (s *Storage) Upload() {

}

func (s *Storage) Delete() {

}
