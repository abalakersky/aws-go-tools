package main

import (
	"fmt"
	"flag"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

func main() {
	t := fmt.Sprintf("%v", time.Now().Unix())
	bucket := flag.String("bucket", t, "Bucket Name to list objects from")
	region := flag.String("region", "us-east-1", "Region to connect to.")
	creds := flag.String("creds", "default", "Credentials Profile to use")
	flag.Parse()
	if *bucket == t {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}
	var wg sync.WaitGroup
	keysCh := make(chan string, 10)

	svc := s3.New(session.New(&aws.Config{
		Region: region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	params := &s3.ListObjectsInput{
		Bucket: bucket,
	}
	wg.Add(1)
	go func(param *s3.ListObjectsInput) {
		defer wg.Done()

		err := svc.ListObjectsPages(params,
			func(page *s3.ListObjectsOutput, last bool) bool {
				for _, object := range page.Contents {
					//keysCh <- fmt.Sprintf("%s:%s", *params.Bucket, *object.Key)
					keysCh <- fmt.Sprintf("%s", *object.Key)
				}
				return true
			},
		)
		if err != nil {
			fmt.Println("Error listing", *bucket, "Objects:", err)
		}
	}(params)
	go func() {
		wg.Wait()
		close(keysCh)
	}()
	for key := range keysCh {
		fmt.Println(key)
	}
}
