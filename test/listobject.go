package main

import (
	"fmt"
	"sync"

	"flag"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	bucket = flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	region = flag.String("region", "us-east-1", "Region to connect to.")
	creds  = flag.String("creds", "default", "Credentials Profile to use")
	search = flag.String("search", "", "Search string to find in object paths")
)

func caseInsesitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	flag.Parse()
	svc := s3.New(session.New(&aws.Config{
		Region:      region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	//	bucket := "pso-training"
	numWorkers := 5

	prefixCh := make(chan string, numWorkers)
	objCh := make(chan *s3.Object, 100)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		// Spin up each worker
		wg.Add(1)
		go func() {
			listObjectsWorker(objCh, prefixCh, *bucket, svc)
			wg.Done()
		}()
	}
	go func() {
		// Wait until workers are finished then close object channel
		wg.Wait()
		close(objCh)
	}()

	go func() {
		if err := getBucketCommonPrefixes(prefixCh, *bucket, svc); err != nil {
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
		params := &s3.ListObjectsInput{
			Bucket:    &bucket,
			Prefix:    &prefix,
			Delimiter: aws.String("/"),
		}
		err := svc.ListObjectsPages(params,
			func(page *s3.ListObjectsOutput, last bool) bool {
				for _, object := range page.Contents {
					objCh <- object
				}
				return true
			},
		)

		//		result, err := svc.ListObjectsPages(&s3.ListObjectsInput{
		//			Bucket: &bucket, Delimiter: aws.String("/"),
		//			Prefix: &prefix,
		//		}),
		if err != nil {
			fmt.Println("failed to list objects by prefix", prefix, err)
			continue
		}
		//		for _, obj := range result.Contents {
		//			objCh <- obj
		//		}
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
