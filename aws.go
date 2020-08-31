package streamcommons

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
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

func awsClient() *s3.Client {
	cfg := awsDefaultConfig.Copy()
	cfg.Region = awsRegion()
	return s3.New(cfg)
}

// GetS3Object gets object from s3 bucket.
func GetS3Object(ctx context.Context, key string) (io.ReadCloser, error) {
	c := awsClient()
	req := c.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsBucket()),
		Key:    &key,
	})
	resp, serr := req.Send(ctx)
	if serr != nil {
		if aerr, ok := serr.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			// Key did not found
			return nil, nil
		}
		return nil, fmt.Errorf("GetS3Object: %v", serr)
	}
	return resp.Body, nil
}

// PutS3Object puts object
func PutS3Object(ctx context.Context, key string, body io.Reader) (err error) {
	cfg := awsDefaultConfig.Copy()
	cfg.Region = awsRegion()
	uploader := s3manager.NewUploader(cfg)
	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(awsBucket()),
		Key:    &key,
		Body:   body,
	})
	return
}

// S3ListV2 downloads object metas which key starts from the given string.
// func S3ListV2(ctx context.Context, prefix string) ([]string, error) {
// 	cfg := awsDefaultConfig.Copy()
// 	cfg.Region = awsRegion()
// 	c := s3.New(cfg)
// 	req := c.ListObjectsV2Request(&s3.ListObjectsV2Input{
// 		Bucket: aws.String(awsBucket()),
// 		Prefix: &prefix,
// 	})
// 	res, serr := req.Send(ctx)
// 	if serr != nil {
// 		return nil, serr
// 	}
// 	contents := res.Contents
// 	keys := make([]string, len(contents))
// 	for i, obj := range contents {
// 		keys[i] = *obj.Key
// 	}
// 	return keys, nil
// }

// S3GetConcurrent is gets multiple objects from S3 in concurrent way.
type S3GetConcurrent struct {
	s3        *s3.Client
	bucketKey string
	keys      []string
	bodies    chan io.ReadCloser
	ctx       context.Context
	cancel    context.CancelFunc
	errc      chan error
	lastErr   error
	closed    bool
}

type s3GetConcurrentResult struct {
	key  string
	body io.ReadCloser
	err  error
}

func (c *S3GetConcurrent) downloadRoutine(key string, result chan s3GetConcurrentResult) {
	req := c.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &c.bucketKey,
		Key:    &key,
	})
	resp, serr := req.Send(c.ctx)
	if serr != nil {
		if aerr, ok := serr.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			// Object did not found
			result <- s3GetConcurrentResult{
				key: key,
			}
			return
		}
		result <- s3GetConcurrentResult{
			key: key,
			err: serr,
		}
		return
	}
	result <- s3GetConcurrentResult{
		key:  key,
		body: resp.Body,
	}
}

func (c *S3GetConcurrent) managerRoutine() {
	var err error
	defer func() {
		close(c.bodies)
		if err != nil {
			c.errc <- fmt.Errorf("managerRoutine: %v", err)
		}
		close(c.errc)
	}()
	// Make a map of index and its object key
	keyMap := make(map[string]int)
	for i, key := range c.keys {
		keyMap[key] = i
	}
	results := make(chan s3GetConcurrentResult)
	running := len(c.keys)
	defer func() {
		// Wait for all child-routine to stop
		for running > 0 {
			r := <-results
			running--
			// Ignore errors
			if r.body != nil {
				serr := r.body.Close()
				if serr != nil {
					if err != nil {
						err = fmt.Errorf("defer body: %v, originally: %v", serr, err)
					} else {
						err = fmt.Errorf("defer body: %v", serr)
					}
				}
			}
		}
	}()
	for _, key := range c.keys {
		go c.downloadRoutine(key, results)
	}
	buffer := make([]*s3GetConcurrentResult, len(c.keys))
	// Current position in the buffer
	n := 0
	defer func() {
		// Close all buffered bodies
		for ; n < len(c.keys); n++ {
			if buffer[n] == nil || buffer[n].body == nil {
				continue
			}
			serr := buffer[n].body.Close()
			if serr == nil {
				continue
			}
			if err != nil {
				err = fmt.Errorf("buffer body: %v, originally: %v", serr, err)
			} else {
				err = fmt.Errorf("buffer body: %v", serr)
			}
		}
	}()
	for running > 0 {
		select {
		case r := <-results:
			running--
			if r.err != nil {
				err = r.err
				return
			}
			i := keyMap[r.key]
			buffer[i] = &r
			for ; n < len(c.keys) && buffer[n] != nil; n++ {
				select {
				case c.bodies <- buffer[n].body:
				case <-c.ctx.Done():
					err = c.ctx.Err()
					return
				}
			}
		case <-c.ctx.Done():
			err = c.ctx.Err()
			return
		}
	}
}

// S3GetAll downloads all objects with the given keys concurrently.
func S3GetAll(ctx context.Context, keys []string) *S3GetConcurrent {
	c := new(S3GetConcurrent)
	c.errc = make(chan error)
	c.bodies = make(chan io.ReadCloser)
	c.keys = keys
	c.bucketKey = awsBucket()
	c.ctx, c.cancel = context.WithCancel(ctx)
	// Create new S3 client
	cfg := awsDefaultConfig.Copy()
	cfg.Region = awsRegion()
	c.s3 = s3.New(cfg)
	go c.managerRoutine()
	return c
}

// Next returns next reader for body.
// Returns nil if object is exhausted.
func (c *S3GetConcurrent) Next() (body io.ReadCloser, ok bool) {
	body, ok = <-c.bodies
	return
}

// Close frees resources associated with this downloader.
func (c *S3GetConcurrent) Close() error {
	if c.closed {
		return c.lastErr
	}
	c.cancel()
	serr, ok := <-c.errc
	if !ok {
		c.lastErr = serr
	}
	return serr
}

func init() {
	// Setup
	var serr error
	awsDefaultConfig, serr = external.LoadDefaultAWSConfig()
	if serr != nil {
		panic(serr)
	}
}
