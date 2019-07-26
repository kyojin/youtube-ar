package storage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Storage is the storage interface.
type Storage interface {
	Upload(context.Context, string) (string, error)
}

// New returns a new storage.
func New() (Storage, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &storage{
		bucket: os.Getenv("S3_BUCKET"),
		s3:     s3.New(s),
	}, nil
}

type storage struct {
	bucket string
	s3     *s3.S3
}

func (s *storage) Upload(ctx context.Context, path string) (string, error) {
	// TODO: add logs

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	input := &s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filepath.Base(path)),
	}
	if _, err := s.s3.PutObjectWithContext(ctx, input); err != nil {
		return "", err
	}
	return path, nil
}
