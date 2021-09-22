package storage

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

const (
	NewsBucket    = "news"
	ContentBucket = "content"
	LogoBucket    = "logo"
	LogoFolder    = "logo"
)

type Config struct {
	Endpoint       string
	AccessKeyId    string
	SecretAccesKey string
	UseSSL         bool
}

type Storage struct {
	Client *minio.Client
}

var Location string = "ap-south-1"

func NewStorage(cfg Config) (*Storage, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.SecretAccesKey, ""),
		Secure: cfg.UseSSL,
	})
	return &Storage{Client: minioClient}, err
}

func (s *Storage) MakeBucket(ctx context.Context, bucketName string) error {
	logrus.Printf("location: %s", Location)
	err := s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: Location, ObjectLocking: false})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			logrus.Printf("We already own %s", bucketName)
		} else {
			return err
		}
	} else {
		logrus.Printf("Successfully created %s", bucketName)
	}
	return nil
}

func (s *Storage) UploadAuthorsLogo(ctx context.Context, bucketName string) error {
	for i := 1; i <= 5; i++ {
		filename := strconv.Itoa(i) + ".jpg"
		path := LogoFolder + "/" + filename
		imageReader, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = s.Client.PutObject(ctx, bucketName, filename, imageReader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			return err
		}
		logrus.Printf("Successfully upload authors logo")
	}
	return nil
}

func getImageReader(URL string) (io.Reader, error) {
	if resp, err := http.Get(URL); err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}
