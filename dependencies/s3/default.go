package s3

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/will-lol/personalWebsiteAwesome/services/blog"
)

type S3 struct {
	bucketName *string
	client     *s3.Client
}

func NewS3() (s *S3, err error) {
	bucketName, err := getBucketName()
	if err != nil {
		return nil, err
	}

	s = &S3{
		bucketName: &bucketName,
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s.client = s3.NewFromConfig(awsConfig)

	return s, nil
}

func getBucketName() (name string, err error) {
	name = os.Getenv("BLOG_BUCKET_NAME")
	if name == "" {
		err = errors.New("Environment variable unset")
	}
	return name, err
}

func (s S3) GetAllFiles() (*[]*blog.SimpleFile, error) {
	slog.Default().Debug("get awwwl the files")
	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: s.bucketName,
	})
	slog.Default().Debug(*res.Contents[0].Key)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	objCh := make(chan *blog.SimpleFile, len(res.Contents))
	errCh := make(chan error, len(res.Contents))

	for _, object := range res.Contents {
		wg.Add(1)

		go func(key *string, errCh chan error, objCh chan *blog.SimpleFile) {
			defer wg.Done()

			bytes, err := s.getObject(key)
			if err != nil {
				errCh <- err
			}

			obj := &blog.SimpleFile{
				Bytes: bytes,
				Name: *key,
			}
			objCh <- obj

			slog.Default().Debug(string(bytes))
		}(object.Key, errCh, objCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(objCh)
	}()

	for err := range errCh {
		return nil, err
	}

	p := make([]*blog.SimpleFile, 0, len(res.Contents))
	for obj := range objCh {
		slog.Default().Debug(string(obj.Bytes))
		p = append(p, obj)
	}

	wg.Wait()

	return &p, nil
}

func (s S3) getObject(key *string) ([]byte, error) {
	slog.Default().Debug("welcome to getObject")
	res, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: s.bucketName,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
