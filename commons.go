package streamcommons

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"unsafe"

	"github.com/aws/aws-lambda-go/events"
)

// MakeResponse generates response with statusCode as HTTP status code and body
func MakeResponse(statusCode int, body string) *events.APIGatewayProxyResponse {
	headers := make(map[string]string)

	headers["Content-Type"] = "text/plain"

	return &events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       body,
		StatusCode: statusCode,
	}
}

// MakeLargeResponse makes response struct with large body, might do compression
func MakeLargeResponse(statusCode int, body []byte, quotaUsed int64) (response *events.APIGatewayProxyResponse, err error) {
	headers := make(map[string]string)

	headers["Content-Type"] = "text/plain"
	headers["ExcDataset-Quota-Used"] = fmt.Sprintf("%d", quotaUsed)

	// if body size is bigger than 5MB, compress data
	if len(body) >= 5*1024*1024 {
		headers["Content-Encoding"] = "gzip"

		buf := make([]byte, 0, 10*1024*1024)
		buffer := bytes.NewBuffer(buf)
		bwriter := base64.NewEncoder(base64.StdEncoding, buffer)
		defer func() {
			cerr := bwriter.Close()
			if cerr != nil {
				if err != nil {
					err = fmt.Errorf("failed to close bwriter, original error was: %s", err.Error())
				} else {
					err = errors.New("failed to close bwriter")
				}
			}
		}()
		// we need to compress return data to clear the limit of 6MB return size
		gwriter, _ := gzip.NewWriterLevel(bwriter, gzip.DefaultCompression)
		defer func() {
			// if gwriter is already closed, then it will return nil as error and do nothing
			cerr := gwriter.Close()
			if cerr != nil {
				if err != nil {
					err = fmt.Errorf("failed to close gwriter, original error was: %s", err.Error())
				} else {
					err = errors.New("failed to close gwriter")
				}
			}
		}()

		gwriter.Write(body)

		ferr := gwriter.Flush()
		if ferr != nil {
			if err != nil {
				err = fmt.Errorf("failed to flush gwriter, original error was: %s", err.Error())
			} else {
				err = errors.New("failed to flush gwriter")
			}
			return
		}
		cerr := gwriter.Close()
		if cerr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close gwriter, original error was: %s", err.Error())
			} else {
				err = errors.New("failed to close gwriter")
			}
			return
		}
		cerr = bwriter.Close()
		if cerr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close bwriter, original error was: %s", err.Error())
			} else {
				err = errors.New("failed to close bwriter")
			}
			return
		}

		return &events.APIGatewayProxyResponse{
			Headers:         headers,
			Body:            buffer.String(),
			StatusCode:      statusCode,
			IsBase64Encoded: true,
		}, nil
	}

	// no data compression
	return &events.APIGatewayProxyResponse{
		Headers:         headers,
		Body:            *(*string)(unsafe.Pointer(&body)),
		StatusCode:      statusCode,
		IsBase64Encoded: false,
	}, nil
}
