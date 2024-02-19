package s3

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/will-lol/personalWebsiteAwesome/services/blog"
)

type s3 struct {
	bucketName *string
	client     *awsS3.Client
}

func NewS3() (s *s3, err error) {
	bucketName, err := getBucketName()
	if err != nil {
		return nil, err
	}

	s = &s3{
		bucketName: &bucketName,
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s.client = awsS3.NewFromConfig(awsConfig)

	return s, nil
}

func getBucketName() (name string, err error) {
	name = os.Getenv("BLOG_BUCKET_NAME")
	if name == "" {
		err = errors.New("Environment variable unset")
	}
	return name, err
}

func (s s3) GetAllFiles() (*[]*blog.SimpleFile, error) {
	res, err := s.client.ListObjectsV2(context.TODO(), &awsS3.ListObjectsV2Input{
		Bucket: s.bucketName,
	})
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
		p = append(p, obj)
	}

	wg.Wait()

	return &p, nil
}

func (s s3) getObject(key *string) ([]byte, error) {
	res, err := s.client.GetObject(context.TODO(), &awsS3.GetObjectInput{
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
