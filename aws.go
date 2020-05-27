package streamcommons

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var cfg aws.Config
var s3serv *s3.Client

// GetS3Object provides new reader for specified s3 object
func GetS3Object(key string) (io.ReadCloser, error) {
	req := s3serv.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("exchangedataset-data"),
		Key:    aws.String(key),
	})
	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Sending GetObject request to s3 failed: %s", err.Error())
	}
	return resp.Body, nil
}

// PutS3Object puts object
func PutS3Object(key string, body io.Reader) (err error) {
	uploader := s3manager.NewUploaderWithClient(s3serv)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("exchangedataset-data"),
		Key:    aws.String(key),
		Body:   body,
	})
	return
}

func init() {
	// aws setup
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Loading config failed: " + err.Error())
	}
	cfg.Region = *aws.String("us-east-2")
	// create new s3 client
	s3serv = s3.New(cfg)
}
