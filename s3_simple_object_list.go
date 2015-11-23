package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/credentials"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	t := time.Now()
	bucket := flag.String("bucket", "test", "Bucket Name to list objects from")
	region := flag.String("region", "us-east-1", "Region to connect to.")
	creds := flag.String("creds", "default", "Credentials Profile to use")
	flag.Parse()
	//	fmt.Printf("Your bucket name is %q", *bucket)

	svc := s3.New(session.New(&aws.Config{
		Region:      region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	result, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucket})
	if err != nil {
		log.Println("Failed to list objects", err)
		return
	}

	fmt.Printf("Here are the objects in %q bucket on %s\n\n", *bucket, t.Format(time.RFC1123))
	//	fmt.Println(result)
	for _, object := range result.Contents {
		fmt.Printf("%s\n", *object.Key)
	}
}
