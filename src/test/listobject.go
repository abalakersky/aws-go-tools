package main

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func main() {
	svc := s3.New(session.New())

	bucket := "mybucket"
	numWorkers := 5

	prefixCh := make(chan string, numWorkers)
	objCh := make(chan *s3.Object, 100)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		// Spin up each worker
		wg.Add(1)
		go func() {
			listObjectsWorker(objCh, prefixCh, bucket, svc)
			wg.Done()
		}()
	}
	go func() {
		// Wait until workers are finished then close object channel
		wg.Wait()
		close(objCh)
	}()

	go func() {
		if err := getBucketCommonPrefixes(prefixCh, bucket, svc); err != nil {
			fmt.Println("error getting bucket common prefixes", err)
		}
		// Close prefixCh so workers will know when to stop
		close(prefixCh)
	}()

	for obj := range objCh {
		// TODO Do something with the objects.
		fmt.Println("Object:", *obj.Key)
	}
}

func listObjectsWorker(objCh chan<- *s3.Object, prefixCh <-chan string, bucket string, svc s3iface.S3API) {
	for prefix := range prefixCh {
		result, err := svc.ListObjects(&s3.ListObjectsInput{
			Bucket: &bucket, Delimiter: aws.String("/"),
			Prefix: &prefix,
		})
		if err != nil {
			fmt.Println("failed to list objects by prefix", prefix, err)
			continue
		}
		for _, obj := range result.Contents {
			objCh <- obj
		}
	}
}

func getBucketCommonPrefixes(prefixCh chan<- string, bucket string, svc s3iface.S3API) error {
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucket, Delimiter: aws.String("/"),
	})
	if err != nil {
		return err
	}

	for _, commonPrefix := range result.CommonPrefixes {
		prefixCh <- *commonPrefix.Prefix
	}

	return nil
}
