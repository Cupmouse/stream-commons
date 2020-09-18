package streamcommons

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/aws/aws-lambda-go/events"
)

// BinanceDecomposeChannel decomposes a given channel name into a symbol name, stream name.
// If the given channel name does not have either one, it returns error.
func BinanceDecomposeChannel(channel string) (symbol string, stream string, err error) {
	index := strings.IndexRune(channel, '@')
	if index == -1 || index == len(channel)-1 {
		err = errors.New("BinanceDecomposeChannel: channel does not have stream name")
		return
	}
	return channel[:index], channel[index+1:], nil
}

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
func MakeLargeResponse(statusCode int, body []byte, quotaUsed int) (response *events.APIGatewayProxyResponse, err error) {
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
					err = fmt.Errorf("close bwriter, originally: %v", err)
				} else {
					err = errors.New("close bwriter")
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
					err = fmt.Errorf("close gwriter: %v, originally: %v", cerr, err)
				} else {
					err = fmt.Errorf("close gwriter: %v", cerr)
				}
			}
		}()

		gwriter.Write(body)

		ferr := gwriter.Flush()
		if ferr != nil {
			if err != nil {
				err = fmt.Errorf("flush gwriter: %v, originally: %v", ferr, err)
			} else {
				err = fmt.Errorf("flush gwriter: %v", ferr)
			}
			return
		}
		cerr := gwriter.Close()
		if cerr != nil {
			if err != nil {
				err = fmt.Errorf("close gwriter: %v, originally: %v", cerr, err)
			} else {
				err = fmt.Errorf("close gwriter: %v", cerr)
			}
			return
		}
		cerr = bwriter.Close()
		if cerr != nil {
			if err != nil {
				err = fmt.Errorf("close bwriter: %v, originally: %v", cerr, err)
			} else {
				err = fmt.Errorf("close bwriter: %v", cerr)
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
