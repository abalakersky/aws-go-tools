package main

import (
	"flag"
	"fmt"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/credentials"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
	"sync"
	"bufio"
	"os"
	"log"
	"path/filepath"
)

func main() {
	bucket := flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	region := flag.String("region", "us-east-1", "Region to connect to.")
	creds := flag.String("creds", "default", "Credentials Profile to use")
	flag.Parse()
	if *bucket == "" {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	f, err := os.Create(dir + "/manifester.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var wg sync.WaitGroup
	keysCh := make(chan string, 10)

	svc := s3.New(session.New(&aws.Config{
		Region:      region,
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
	w := bufio.NewWriter(f)
	for key := range keysCh {
		fmt.Println(key)
		w.WriteString(key + "\n")
	}
	w.Flush()
}
