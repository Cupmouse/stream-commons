package streamcommons

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

const awsS3TestRegion = "ap-northeast-1"
const awsS3TestBucket = "exchangedataset-test-data"
const awsS3ProductionRegion = "us-east-2"
const awsS3ProductionBucket = "exchangedataset-data"

var awsTest = true
var awsDefaultConfig aws.Config

// AWSEnableProduction enables production mode.
func AWSEnableProduction() {
	awsTest = false
}

func awsRegion() string {
	if awsTest {
		return awsS3TestRegion
	}
	return awsS3ProductionRegion
}

func awsBucket() string {
	if awsTest {
		return awsS3TestBucket
	}
	return awsS3ProductionBucket
}

// S3GetConcurrent is gets multiple objects from S3 in concurrent way.
type S3GetConcurrent struct {
	s3        *s3.Client
	bucketKey string
	keys      []string
	bodies    chan io.Closer
	stop      chan struct{}
	errc      chan error
	lastErr   error
	closed    bool
}

type s3GetConcurrentResult struct {
	key  string
	body io.Closer
	err  error
}

func (c *S3GetConcurrent) downloadRoutine(ctx context.Context, key string, result chan s3GetConcurrentResult) {
	req := c.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &c.bucketKey,
		Key:    &key,
	})
	resp, serr := req.Send(ctx)
	if serr != nil {
		result <- s3GetConcurrentResult{
			err: serr,
		}
		return
	}
	result <- s3GetConcurrentResult{
		key:  key,
		body: resp.Body,
	}
}

func (c *S3GetConcurrent) managerRoutine(ctx context.Context) {
	var err error
	defer func() {
		close(c.bodies)
		if err != nil {
			c.errc <- err
		}
		close(c.errc)
	}()
	// Make a map of index and its object key
	keyMap := make(map[string]int)
	for i, key := range c.keys {
		keyMap[key] = i
	}
	ctxChild, cancel := context.WithCancel(ctx)
	defer cancel()
	results := make(chan s3GetConcurrentResult)
	for _, key := range c.keys {
		go c.downloadRoutine(ctxChild, key, results)
	}
	buffer := make([]io.Closer, len(c.keys))
	left := len(c.keys)
	// Next index
	n := 0
	defer func() {
		for left > 0 {
			r := <-results
			if r.body != nil {
				serr := r.body.Close()
				if serr != nil {
					if err != nil {
						err = fmt.Errorf("%v, originally: %v", serr, err)
					} else {
						err = serr
					}
				}
			}
			left--
		}
	}()
	for left > 0 {
		select {
		case r := <-results:
			if r.err != nil {
				err = r.err
				return
			}
			i := keyMap[r.key]
			buffer[i] = r.body
			for ; buffer[n] != nil && n < len(c.keys); n++ {
				c.bodies <- buffer[n]
			}
			c.bodies <- r.body
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
		left--
	}
}

// S3GetAll downloads all objects with the given keys concurrently.
func S3GetAll(ctx context.Context, keys []string) *S3GetConcurrent {
	c := new(S3GetConcurrent)
	c.errc = make(chan error)
	c.bodies = make(chan io.Closer)
	c.keys = keys
	c.bucketKey = awsBucket()
	// Create new S3 client
	cfg := awsDefaultConfig.Copy()
	cfg.Region = awsRegion()
	c.s3 = s3.New(cfg)
	go c.managerRoutine(ctx)
	return c
}

// Next returns next reader for body.
// Returns nil if object is exhausted.
func (c *S3GetConcurrent) Next() io.Closer {
	return <-c.bodies
}

// Close frees resources associated with this downloader.
func (c *S3GetConcurrent) Close() error {
	if c.closed {
		return c.lastErr
	}
	close(c.stop)
	serr, ok := <-c.errc
	if !ok {
		c.lastErr = serr
	}
	return serr
}

// PutS3Object puts object
func PutS3Object(key string, body io.Reader) (err error) {
	cfg := awsDefaultConfig.Copy()
	cfg.Region = awsRegion()
	uploader := s3manager.NewUploader(cfg)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("exchangedataset-data"),
		Key:    aws.String(key),
		Body:   body,
	})
	return
}

func init() {
	// Setup
	var serr error
	awsDefaultConfig, serr = external.LoadDefaultAWSConfig()
	if serr != nil {
		panic(serr)
	}
}
